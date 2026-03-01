package com.test.user.controller;

import com.test.user.service.InternalShowService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/internal/shows")
@RequiredArgsConstructor
@Slf4j
public class InternalShowController {

    private final InternalShowService internalShowService;

    @GetMapping("/{showId}/booked-seats")
    public ResponseEntity<List<String>> getBookedSeats(@PathVariable Long showId) {
        log.info("Internal API called: Get booked seats for showId: {}", showId);
        List<String> bookedSeats = internalShowService.getConfirmedBookedSeats(showId);
        return ResponseEntity.ok(bookedSeats);
    }
}
