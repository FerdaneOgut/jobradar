using Microsoft.AspNetCore.Mvc;
using UserService.DTOs;
using UserService.Services;

namespace UserService.Controllers;

[ApiController]
[Route("api/[controller]")]
public class AuthController : ControllerBase
{
    private readonly UserAuthService _userService;

    public AuthController(UserAuthService userService)
    {
        _userService = userService;
    }

    [HttpPost("register")]
    public async Task<IActionResult> Register([FromBody] RegisterRequest request)
    {
        var existing = await _userService.GetByEmailAsync(request.Email);
        if (existing != null)
            return BadRequest(new { message = "Email already exists" });

        var user = await _userService.CreateAsync(request);
        var token = _userService.GenerateJwtToken(user);

        return Ok(new AuthResponse
        {
            Token = token,
            FullName = user.FullName,
            Email = user.Email,
            Plan = user.Plan
        });
    }

    [HttpPost("login")]
    public async Task<IActionResult> Login([FromBody] LoginRequest request)
    {
        var user = await _userService.GetByEmailAsync(request.Email);
        if (user == null)
            return Unauthorized(new { message = "Invalid credentials" });

        if (!_userService.VerifyPassword(request.Password, user.PasswordHash))
            return Unauthorized(new { message = "Invalid credentials" });

        var token = _userService.GenerateJwtToken(user);

        return Ok(new AuthResponse
        {
            Token = token,
            FullName = user.FullName,
            Email = user.Email,
            Plan = user.Plan
        });
    }
}