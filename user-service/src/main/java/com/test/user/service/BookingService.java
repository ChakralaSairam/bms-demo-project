package com.test.user.service;

import com.test.user.client.InventoryClient;
import com.test.user.client.SeatLockRequest;
import com.test.user.client.SeatLockResponse;
import com.test.user.client.SeatUnlockResponse;
import com.test.user.dto.BookingRequest;
import com.test.user.dto.BookingResponse;
import com.test.user.dto.CancelBookingRequest;
import com.test.user.dto.PaymentRequest;
import com.test.user.dto.PaymentResponse;
import com.test.user.exception.ShowNotFoundException;
import com.test.user.models.Booking;
import com.test.user.models.BookingStatus;
import com.test.user.models.Show;
import com.test.user.models.User;
import com.test.user.repository.BookingRepository;
import com.test.user.repository.ShowRepository;
import com.test.user.repository.UserRepository;
import jakarta.servlet.http.HttpServletRequest;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.context.request.RequestContextHolder;
import org.springframework.web.context.request.ServletRequestAttributes;

@Service
@RequiredArgsConstructor
@Slf4j
public class BookingService {

    private final BookingRepository bookingRepository;
    private final ShowRepository showRepository;
    private final UserRepository userRepository;
    private final InventoryClient inventoryClient;

    @Transactional
    public BookingResponse initiateBooking(BookingRequest request) {
        log.info("Initiating booking for showId: {}, seatId: {}", request.getShowId(), request.getSeatId());
        
        Show show = showRepository.findById(request.getShowId())
                .orElseThrow(() -> new ShowNotFoundException("Show not found with ID: " + request.getShowId()));

        Long userId = getCurrentUserId();
        String token = getAuthorizationHeader();
        log.info("User ID: {}, Token: {}", userId, token != null ? token.substring(0, 20) + "..." : "null");
        
        SeatLockRequest lockRequest = SeatLockRequest.builder()
                .showId(request.getShowId())
                .seatId(request.getSeatId())
                .userId(userId)
                .build();
        
        log.info("Calling inventory service - URL: http://192.168.1.3:8081/inventory/lock");
        log.info("Request payload - showId: {}, seatId: {}, userId: {}", 
                lockRequest.getShowId(), lockRequest.getSeatId(), lockRequest.getUserId());
        
        try {
            SeatLockResponse lockResponse = inventoryClient.lockSeat(token, lockRequest);
            log.info("✅ Inventory service responded successfully");
            log.info("Response status: {}", lockResponse.getStatus());
            if (lockResponse.getShowId() != null) {
                log.info("Show: {}, Seat: {}, Held for: {}", 
                        lockResponse.getShowId(), lockResponse.getSeatId(), lockResponse.getHeldFor());
            }
        } catch (Exception e) {
            log.error("❌ Failed to lock seat - Error: {}", e.getMessage(), e);
            throw e;
        }

        Booking booking = Booking.builder()
                .userId(userId)
                .showId(request.getShowId())
                .seatId(request.getSeatId())
                .status(BookingStatus.PENDING)
                .amount(show.getPrice())
                .build();

        booking = bookingRepository.save(booking);
        log.info("Booking created successfully with ID: {}", booking.getId());

        return BookingResponse.builder()
                .bookingId(booking.getId())
                .status("PENDING")
                .message("Booking initiated successfully")
                .build();
    }

    @Transactional
    public PaymentResponse processPayment(PaymentRequest request) {
        log.info("Processing payment for bookingId: {}", request.getBookingId());
        
        Booking booking = bookingRepository.findById(request.getBookingId())
                .orElseThrow(() -> new RuntimeException("Booking not found with ID: " + request.getBookingId()));
        
        if (booking.getStatus() != BookingStatus.PENDING) {
            throw new RuntimeException("Booking is not in PENDING status");
        }
        
        // Simulate payment processing
        log.info("Simulating payment for amount: {}", booking.getAmount());
        boolean paymentSuccess = true; // Simulated success
        
        if (paymentSuccess) {
            booking.setStatus(BookingStatus.CONFIRMED);
            bookingRepository.save(booking);
            log.info("Payment successful, booking confirmed");
            
            // Unlock seat in inventory service
            String token = getAuthorizationHeader();
            try {
                inventoryClient.unlockSeat(token, booking.getShowId(), booking.getSeatId());
                log.info("Seat unlocked in inventory service");
            } catch (Exception e) {
                log.error("Failed to unlock seat: {}", e.getMessage());
            }
            
            return PaymentResponse.builder()
                    .bookingId(booking.getId())
                    .status("CONFIRMED")
                    .amount(booking.getAmount())
                    .message("Payment successful")
                    .build();
        } else {
            booking.setStatus(BookingStatus.FAILED);
            bookingRepository.save(booking);
            log.info("Payment failed");
            
            return PaymentResponse.builder()
                    .bookingId(booking.getId())
                    .status("FAILED")
                    .amount(booking.getAmount())
                    .message("Payment failed")
                    .build();
        }
    }

    private Long getCurrentUserId() {
        String username = SecurityContextHolder.getContext().getAuthentication().getName();
        User user = userRepository.findByUsername(username)
                .orElseThrow(() -> new RuntimeException("User not found"));
        return user.getId();
    }

    private String getAuthorizationHeader() {
        ServletRequestAttributes attributes = (ServletRequestAttributes) RequestContextHolder.getRequestAttributes();
        if (attributes != null) {
            HttpServletRequest request = attributes.getRequest();
            String authHeader = request.getHeader("Authorization");
            if (authHeader != null && authHeader.startsWith("Bearer ")) {
                log.debug("Found Authorization header in request");
                return authHeader;
            }
        }
        log.warn("No Authorization header found in request");
        return "Bearer ";
    }

    @Transactional
    public String cancelBooking(CancelBookingRequest request) {
        log.info("Cancelling booking for showId: {}, seatId: {}", request.getShowId(), request.getSeatId());
        
        String token = getAuthorizationHeader();
        
        SeatUnlockResponse unlockResponse = inventoryClient.unlockSeat(
                token, 
                request.getShowId(), 
                request.getSeatId()
        );
        
        log.info("Seat lock released - Show: {}, Seat: {}", 
                unlockResponse.getShowId(), unlockResponse.getSeatId());
        
        return "Booking cancelled successfully";
    }
}
