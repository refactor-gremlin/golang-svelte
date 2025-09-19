using Microsoft.AspNetCore.Mvc;

namespace MySvelteApp.Server.Presentation.Controllers;

[ApiController]
[Route("[controller]")]
public class TestAuthController : ControllerBase
{
    [HttpGet]
    public IActionResult Get()
    {
        return Ok(new { Message = "If you can see this, you are authenticated!" });
    }
}
