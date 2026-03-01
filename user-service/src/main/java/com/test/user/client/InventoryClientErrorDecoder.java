package com.test.user.client;

import com.test.user.exception.InventoryServiceException;
import com.test.user.exception.SeatAlreadyLockedException;
import feign.Response;
import feign.codec.ErrorDecoder;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

@Component
@Slf4j
public class InventoryClientErrorDecoder implements ErrorDecoder {

    private final ErrorDecoder defaultDecoder = new Default();

    @Override
    public Exception decode(String methodKey, Response response) {
        log.error("Inventory service error - Status: {}, Method: {}", response.status(), methodKey);
        
        return switch (response.status()) {
            case 404 -> {
                log.warn("Resource not found - Status 404");
                yield new InventoryServiceException("Seat lock not found or already expired");
            }
            case 409 -> {
                log.warn("Seat already locked - Status 409");
                yield new SeatAlreadyLockedException("Seat is already locked by another user");
            }
            case 400 -> {
                log.error("Bad request to inventory service - Status 400");
                yield new InventoryServiceException("Invalid request to inventory service");
            }
            case 401 -> {
                log.error("Unauthorized access to inventory service - Status 401");
                yield new InventoryServiceException("Authentication failed with inventory service");
            }
            case 500, 503 -> {
                log.error("Inventory service unavailable - Status {}", response.status());
                yield new InventoryServiceException("Inventory service is unavailable");
            }
            default -> defaultDecoder.decode(methodKey, response);
        };
    }
}
