package com.test.user.controller;

import com.test.user.dto.TheaterDTO;
import com.test.user.service.TheaterAdminService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/admin/theaters")
@RequiredArgsConstructor
public class TheaterAdminController {

    private final TheaterAdminService theaterAdminService;

    @GetMapping
    @PreAuthorize("hasAuthority('ROLE_ADMIN')")
    public ResponseEntity<List<TheaterDTO>> getAllTheaters() {
        return ResponseEntity.ok(theaterAdminService.getAllTheaters());
    }

    @PostMapping
    @PreAuthorize("hasAuthority('ROLE_ADMIN')")
    public ResponseEntity<TheaterDTO> createTheater(@Valid @RequestBody TheaterDTO dto) {
        return ResponseEntity.status(HttpStatus.CREATED).body(theaterAdminService.createTheater(dto));
    }

    @PutMapping("/{id}")
    @PreAuthorize("hasAuthority('ROLE_ADMIN')")
    public ResponseEntity<TheaterDTO> updateTheater(@PathVariable Long id, @Valid @RequestBody TheaterDTO dto) {
        return ResponseEntity.ok(theaterAdminService.updateTheater(id, dto));
    }

    @DeleteMapping("/{id}")
    @PreAuthorize("hasAuthority('ROLE_ADMIN')")
    public ResponseEntity<Void> deleteTheater(@PathVariable Long id) {
        theaterAdminService.deleteTheater(id);
        return ResponseEntity.noContent().build();
    }
}
