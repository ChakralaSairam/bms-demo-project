package com.test.user.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.*;

@Getter @Setter
@NoArgsConstructor
@AllArgsConstructor
public class BookingRequest {
    @NotNull(message = "Show ID is required")
    private Long showId;
    
    @NotBlank(message = "Seat ID is required")
    private String seatId;
}
