package com.test.user.service;

import com.test.user.dto.TheaterDTO;
import com.test.user.models.Theater;
import com.test.user.repository.TheaterRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
@Slf4j
public class TheaterAdminService {

    private final TheaterRepository theaterRepository;

    @Transactional(readOnly = true)
    public List<TheaterDTO> getAllTheaters() {
        return theaterRepository.findAll().stream()
                .map(this::toDTO)
                .collect(Collectors.toList());
    }

    @Transactional
    public TheaterDTO createTheater(TheaterDTO dto) {
        log.info("Creating theater: {}", dto.getName());
        Theater theater = Theater.builder()
                .name(dto.getName())
                .city(dto.getCity())
                .build();
        theater = theaterRepository.save(theater);
        return toDTO(theater);
    }

    @Transactional
    public TheaterDTO updateTheater(Long id, TheaterDTO dto) {
        log.info("Updating theater ID: {}", id);
        Theater theater = theaterRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Theater not found"));
        theater.setName(dto.getName());
        theater.setCity(dto.getCity());
        theater = theaterRepository.save(theater);
        return toDTO(theater);
    }

    @Transactional
    public void deleteTheater(Long id) {
        log.info("Deleting theater ID: {}", id);
        theaterRepository.deleteById(id);
    }

    private TheaterDTO toDTO(Theater theater) {
        return TheaterDTO.builder()
                .id(theater.getId())
                .name(theater.getName())
                .city(theater.getCity())
                .build();
    }
}
