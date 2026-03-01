package com.test.user.controller;

import com.test.user.dto.MovieDTO;
import com.test.user.service.MovieAdminService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/admin/movies")
@RequiredArgsConstructor
public class MovieAdminController {

    private final MovieAdminService movieAdminService;

    @PostMapping
    @PreAuthorize("hasAuthority('ROLE_ADMIN')")
    public ResponseEntity<MovieDTO> createMovie(@Valid @RequestBody MovieDTO dto) {
        return ResponseEntity.status(HttpStatus.CREATED).body(movieAdminService.createMovie(dto));
    }

    @PutMapping("/{id}")
    @PreAuthorize("hasAuthority('ROLE_ADMIN')")
    public ResponseEntity<MovieDTO> updateMovie(@PathVariable Long id, @Valid @RequestBody MovieDTO dto) {
        return ResponseEntity.ok(movieAdminService.updateMovie(id, dto));
    }

    @DeleteMapping("/{id}")
    @PreAuthorize("hasAuthority('ROLE_ADMIN')")
    public ResponseEntity<Void> deleteMovie(@PathVariable Long id) {
        movieAdminService.deleteMovie(id);
        return ResponseEntity.noContent().build();
    }
}
