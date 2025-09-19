using System.ComponentModel.DataAnnotations;
using System.Text;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Diagnostics;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Options;
using Microsoft.IdentityModel.Tokens;
using Microsoft.OpenApi.Models;
using MySvelteApp.Server.Application.Authentication;
using MySvelteApp.Server.Application.Common.Interfaces;
using MySvelteApp.Server.Application.Pokemon;
using MySvelteApp.Server.Application.Weather;
using MySvelteApp.Server.Infrastructure.Authentication;
using MySvelteApp.Server.Infrastructure.External;
using MySvelteApp.Server.Infrastructure.Persistence;
using MySvelteApp.Server.Infrastructure.Persistence.Repositories;
using MySvelteApp.Server.Infrastructure.Security;
using MySvelteApp.Server.Infrastructure.Weather;
using OpenTelemetry.Exporter;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using Serilog;
using Serilog.Sinks.Grafana.Loki;

var builder = WebApplication.CreateBuilder(args);

const string WebsiteClientOrigin = "website_client";

builder.Services.AddCors(options =>
    options.AddPolicy(
        WebsiteClientOrigin,
        policy =>
            policy
                .WithOrigins("http://localhost:5173", "http://localhost:3000", "http://web:3000")
                .AllowAnyHeader()
                .AllowAnyMethod()
                .AllowCredentials()
    )
);

// Bind and validate JwtOptions
builder
    .Services.AddOptions<JwtOptions>()
    .Bind(builder.Configuration.GetSection("Jwt"))
    .ValidateDataAnnotations()
    .Validate(o => !string.IsNullOrWhiteSpace(o.Key), "Jwt:Key cannot be blank/whitespace.")
    .Validate(o => !string.IsNullOrWhiteSpace(o.Issuer), "Jwt:Issuer cannot be blank/whitespace.")
    .Validate(
        o => !string.IsNullOrWhiteSpace(o.Audience),
        "Jwt:Audience cannot be blank/whitespace."
    )
    .ValidateOnStart();

builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme).AddJwtBearer();

builder
    .Services.AddOptions<JwtBearerOptions>(JwtBearerDefaults.AuthenticationScheme)
    .Configure<IOptions<JwtOptions>>(
        (options, jwt) =>
        {
            var o = jwt.Value;
            options.TokenValidationParameters = new TokenValidationParameters
            {
                ValidateIssuer = true,
                ValidateAudience = true,
                ValidateLifetime = true,
                ValidateIssuerSigningKey = true,
                ValidIssuer = o.Issuer,
                ValidAudience = o.Audience,
                IssuerSigningKey = new SymmetricSecurityKey(DeriveKeyBytes(o.Key)),
                // Optional hardening: stricter expiry validation
                ClockSkew = TimeSpan.Zero,
            };
        }
    );

// Shared key derivation helper to prevent duplication
static byte[] DeriveKeyBytes(string key)
{
    return key.StartsWith("base64:", StringComparison.Ordinal)
        ? Convert.FromBase64String(key["base64:".Length..])
        : Encoding.UTF8.GetBytes(key);
}

builder
    .Services.AddAuthorizationBuilder()
    .SetFallbackPolicy(new AuthorizationPolicyBuilder().RequireAuthenticatedUser().Build());

builder.Services.AddControllers();
builder.Services.AddProblemDetails();
builder.Services.AddHealthChecks();

builder.Services.AddSwaggerGen(c =>
{
    c.SwaggerDoc("v1", new OpenApiInfo { Title = "MySvelteApp.Server", Version = "v1" });

    c.AddSecurityDefinition(
        "Bearer",
        new OpenApiSecurityScheme
        {
            Name = "Authorization",
            Type = SecuritySchemeType.Http,
            Scheme = "bearer",
            BearerFormat = "JWT",
            In = ParameterLocation.Header,
            Description = "Enter 'Bearer <token>'",
        }
    );

    c.AddSecurityRequirement(
        new OpenApiSecurityRequirement
        {
            {
                new OpenApiSecurityScheme
                {
                    Reference = new OpenApiReference
                    {
                        Type = ReferenceType.SecurityScheme,
                        Id = "Bearer",
                    },
                },
                Array.Empty<string>()
            },
        }
    );

    var xmlFile = $"{System.Reflection.Assembly.GetExecutingAssembly().GetName().Name}.xml";
    var xmlPath = System.IO.Path.Combine(AppContext.BaseDirectory, xmlFile);
    if (System.IO.File.Exists(xmlPath))
    {
        c.IncludeXmlComments(xmlPath);
    }
});

var promtailUrl =
    builder.Configuration["LOKI_PUSH_URL"] ?? "http://localhost:3101/loki/api/v1/push";
var apiServiceName = builder.Configuration["OTEL_SERVICE_NAME"] ?? "mysvelteapp-api";
var environmentName = builder.Environment.EnvironmentName ?? "Development";

builder.Host.UseSerilog(
    (_, _, configuration) =>
        configuration
            .MinimumLevel.Information()
            .Enrich.FromLogContext()
            .Enrich.WithProperty("service", apiServiceName)
            .Enrich.WithProperty("env", environmentName.ToLowerInvariant())
            .WriteTo.Console()
            .WriteTo.GrafanaLoki(promtailUrl)
);

var serviceName = apiServiceName;
var otlpEndpoint =
    builder.Configuration["OTEL_EXPORTER_OTLP_ENDPOINT"] ?? "http://localhost:4318/v1/traces";
var otlpProtocol = builder.Configuration["OTEL_EXPORTER_OTLP_PROTOCOL"] ?? "http/protobuf";

builder
    .Services.AddOpenTelemetry()
    .WithTracing(tracing =>
        tracing
            .SetResourceBuilder(ResourceBuilder.CreateDefault().AddService(serviceName))
            .AddAspNetCoreInstrumentation(options => options.RecordException = true)
            .AddHttpClientInstrumentation()
            .AddOtlpExporter(options =>
            {
                options.Endpoint = new Uri(otlpEndpoint);
                options.Protocol = string.Equals(
                    otlpProtocol,
                    "grpc",
                    StringComparison.OrdinalIgnoreCase
                )
                    ? OtlpExportProtocol.Grpc
                    : OtlpExportProtocol.HttpProtobuf;
            })
    );

builder.Services.AddDbContext<AppDbContext>(options =>
    options.UseInMemoryDatabase("MySvelteAppDb")
);

builder.Services.AddScoped<IJwtTokenGenerator, JwtTokenGenerator>();
builder.Services.AddScoped<IPasswordHasher, PasswordHasher>();
builder.Services.AddScoped<IUserRepository, UserRepository>();
builder.Services.AddScoped<IAuthService, AuthService>();
builder.Services.AddHttpClient<IRandomPokemonService, PokeApiRandomPokemonService>();
builder.Services.AddSingleton<IWeatherForecastService, WeatherForecastService>();

var app = builder.Build();

// Global exception handling with ProblemDetails - must be first middleware
app.UseExceptionHandler(errorApp =>
    errorApp.Run(async context =>
    {
        var exceptionHandler = context.Features.Get<IExceptionHandlerFeature>();
        if (exceptionHandler is not null)
        {
            var problemDetails = new ProblemDetails
            {
                Status = StatusCodes.Status500InternalServerError,
                Title = "An unexpected error occurred.",
                Detail = app.Environment.IsDevelopment() ? exceptionHandler.Error.Message : null,
                Instance = context.Request.Path,
            };

            context.Response.StatusCode = problemDetails.Status.Value;
            context.Response.ContentType = "application/problem+json";
            await context.Response.WriteAsJsonAsync(problemDetails);
        }
    })
);

app.UseCors(WebsiteClientOrigin);

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();
app.UseSerilogRequestLogging();
app.UseAuthentication();
app.UseAuthorization();
app.MapControllers();
app.MapHealthChecks("/health").AllowAnonymous();
app.Run();
