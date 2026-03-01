package com.test.user.service;

import com.test.user.exception.ShowNotFoundException;
import com.test.user.models.BookingStatus;
import com.test.user.repository.BookingRepository;
import com.test.user.repository.ShowRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Service
@RequiredArgsConstructor
@Slf4j
public class InternalShowService {

    private final BookingRepository bookingRepository;
    private final ShowRepository showRepository;

    @Transactional(readOnly = true)
    public List<String> getConfirmedBookedSeats(Long showId) {
        log.info("Fetching confirmed booked seats for showId: {}", showId);
        
        if (!showRepository.existsById(showId)) {
            throw new ShowNotFoundException("Show not found with ID: " + showId);
        }
        
        return bookingRepository.findSeatIdsByShowIdAndStatus(showId, BookingStatus.CONFIRMED);
    }
}
