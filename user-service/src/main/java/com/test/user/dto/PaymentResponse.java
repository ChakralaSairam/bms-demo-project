package com.test.user.dto;

import lombok.*;
import java.math.BigDecimal;
import java.util.UUID;

@Getter @Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class PaymentResponse {
    private UUID bookingId;
    private String status;
    private BigDecimal amount;
    private String message;
}
