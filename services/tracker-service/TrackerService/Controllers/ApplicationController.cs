using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using System.Security.Claims;
using TrackerService.DTOs;
using TrackerService.Services;

namespace TrackerService.Controllers;

[ApiController]
[Route("api/[controller]")]
[Authorize]
public class ApplicationController : ControllerBase
{
    private readonly ApplicationService _applicationService;

    public ApplicationController(ApplicationService applicationService)
    {
        _applicationService = applicationService;
    }

    [HttpPost]
    public async Task<IActionResult> Create([FromBody] CreateApplicationRequest request)
    {
        var userId = User.FindFirst(ClaimTypes.NameIdentifier)?.Value;
        if (userId == null) return Unauthorized();

        var application = await _applicationService.CreateAsync(userId, request);
        return Ok(application);
    }

    [HttpGet]
    public async Task<IActionResult> GetAll()
    {
        var userId = User.FindFirst(ClaimTypes.NameIdentifier)?.Value;
        if (userId == null) return Unauthorized();

        var applications = await _applicationService.GetByUserIdAsync(userId);
        return Ok(applications);
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> GetById(string id)
    {
        var userId = User.FindFirst(ClaimTypes.NameIdentifier)?.Value;
        if (userId == null) return Unauthorized();

        var application = await _applicationService.GetByIdAsync(id);
        if (application == null || application.UserId != userId)
            return NotFound();

        return Ok(application);
    }

    [HttpPut("{id}")]
    public async Task<IActionResult> Update(string id, [FromBody] UpdateApplicationRequest request)
    {
        var userId = User.FindFirst(ClaimTypes.NameIdentifier)?.Value;
        if (userId == null) return Unauthorized();

        var application = await _applicationService.UpdateAsync(id, userId, request);
        if (application == null) return NotFound();

        return Ok(application);
    }

    [HttpDelete("{id}")]
    public async Task<IActionResult> Delete(string id)
    {
        var userId = User.FindFirst(ClaimTypes.NameIdentifier)?.Value;
        if (userId == null) return Unauthorized();

        var deleted = await _applicationService.DeleteAsync(id, userId);
        if (!deleted) return NotFound();

        return Ok(new { message = "Application deleted" });
    }
}