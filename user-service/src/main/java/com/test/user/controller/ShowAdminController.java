package com.test.user.controller;

import com.test.user.dto.ShowDTO;
import com.test.user.service.ShowAdminService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/admin/shows")
@RequiredArgsConstructor
public class ShowAdminController {

    private final ShowAdminService showAdminService;

    @GetMapping
    @PreAuthorize("hasAuthority('ROLE_ADMIN')")
    public ResponseEntity<List<ShowDTO>> getAllShows() {
        return ResponseEntity.ok(showAdminService.getAllShows());
    }

    @PostMapping
    @PreAuthorize("hasAuthority('ROLE_ADMIN')")
    public ResponseEntity<ShowDTO> createShow(@Valid @RequestBody ShowDTO dto) {
        return ResponseEntity.status(HttpStatus.CREATED).body(showAdminService.createShow(dto));
    }

    @PutMapping("/{id}")
    @PreAuthorize("hasAuthority('ROLE_ADMIN')")
    public ResponseEntity<ShowDTO> updateShow(@PathVariable Long id, @Valid @RequestBody ShowDTO dto) {
        return ResponseEntity.ok(showAdminService.updateShow(id, dto));
    }

    @DeleteMapping("/{id}")
    @PreAuthorize("hasAuthority('ROLE_ADMIN')")
    public ResponseEntity<Void> deleteShow(@PathVariable Long id) {
        showAdminService.deleteShow(id);
        return ResponseEntity.noContent().build();
    }
}
