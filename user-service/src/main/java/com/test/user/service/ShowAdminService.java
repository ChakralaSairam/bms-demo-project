package com.test.user.service;

import com.test.user.dto.ShowDTO;
import com.test.user.models.Show;
import com.test.user.repository.ShowRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
@Slf4j
public class ShowAdminService {

    private final ShowRepository showRepository;

    @Transactional(readOnly = true)
    public List<ShowDTO> getAllShows() {
        return showRepository.findAll().stream()
                .map(this::toDTO)
                .collect(Collectors.toList());
    }

    @Transactional
    public ShowDTO createShow(ShowDTO dto) {
        log.info("Creating show for movie ID: {}", dto.getMovieId());
        Show show = Show.builder()
                .movieId(dto.getMovieId())
                .theaterId(dto.getTheaterId())
                .startTime(dto.getStartTime())
                .price(dto.getPrice())
                .build();
        show = showRepository.save(show);
        return toDTO(show);
    }

    @Transactional
    public ShowDTO updateShow(Long id, ShowDTO dto) {
        log.info("Updating show ID: {}", id);
        Show show = showRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Show not found"));
        show.setMovieId(dto.getMovieId());
        show.setTheaterId(dto.getTheaterId());
        show.setStartTime(dto.getStartTime());
        show.setPrice(dto.getPrice());
        show = showRepository.save(show);
        return toDTO(show);
    }

    @Transactional
    public void deleteShow(Long id) {
        log.info("Deleting show ID: {}", id);
        showRepository.deleteById(id);
    }

    private ShowDTO toDTO(Show show) {
        return ShowDTO.builder()
                .id(show.getId())
                .movieId(show.getMovieId())
                .theaterId(show.getTheaterId())
                .startTime(show.getStartTime())
                .price(show.getPrice())
                .build();
    }
}
