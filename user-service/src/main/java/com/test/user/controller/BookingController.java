package com.test.user.controller;

import com.test.user.dto.BookingRequest;
import com.test.user.dto.BookingResponse;
import com.test.user.dto.CancelBookingRequest;
import com.test.user.dto.PaymentRequest;
import com.test.user.dto.PaymentResponse;
import com.test.user.service.BookingService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/bookings")
@RequiredArgsConstructor
public class BookingController {

    private final BookingService bookingService;

    @PostMapping("/initiate")
    public ResponseEntity<BookingResponse> initiateBooking(@Valid @RequestBody BookingRequest request) {
        BookingResponse response = bookingService.initiateBooking(request);
        return ResponseEntity.status(HttpStatus.CREATED).body(response);
    }

    @PostMapping("/payment")
    public ResponseEntity<PaymentResponse> processPayment(@Valid @RequestBody PaymentRequest request) {
        PaymentResponse response = bookingService.processPayment(request);
        return ResponseEntity.ok(response);
    }

    @DeleteMapping("/cancel")
    public ResponseEntity<String> cancelBooking(@Valid @RequestBody CancelBookingRequest request) {
        String message = bookingService.cancelBooking(request);
        return ResponseEntity.ok(message);
    }
}
