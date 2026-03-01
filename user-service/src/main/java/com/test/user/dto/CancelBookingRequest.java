package com.test.user.dto;

import jakarta.validation.constraints.NotNull;
import lombok.*;

@Getter @Setter
@NoArgsConstructor
@AllArgsConstructor
public class CancelBookingRequest {
    @NotNull(message = "Show ID is required")
    private Long showId;
    
    @NotNull(message = "Seat ID is required")
    private String seatId;
}
