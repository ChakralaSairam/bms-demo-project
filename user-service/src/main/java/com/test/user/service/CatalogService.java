package com.test.user.service;

import com.test.user.dto.MovieDTO;
import com.test.user.dto.ShowDTO;
import com.test.user.models.Movie;
import com.test.user.models.Show;
import com.test.user.repository.MovieRepository;
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
public class CatalogService {

    private final MovieRepository movieRepository;
    private final ShowRepository showRepository;

    @Transactional(readOnly = true)
    public List<MovieDTO> getAllMovies() {
        log.info("Fetching all movies");
        return movieRepository.findAll().stream()
                .map(this::toMovieDTO)
                .collect(Collectors.toList());
    }

    @Transactional(readOnly = true)
    public List<ShowDTO> getShowsByMovieId(Long movieId) {
        log.info("Fetching shows for movie ID: {}", movieId);
        return showRepository.findByMovieId(movieId).stream()
                .map(this::toShowDTO)
                .collect(Collectors.toList());
    }

    private MovieDTO toMovieDTO(Movie movie) {
        return MovieDTO.builder()
                .id(movie.getId())
                .title(movie.getTitle())
                .genre(movie.getGenre())
                .duration(movie.getDuration())
                .language(movie.getLanguage())
                .build();
    }

    private ShowDTO toShowDTO(Show show) {
        return ShowDTO.builder()
                .id(show.getId())
                .movieId(show.getMovieId())
                .theaterId(show.getTheaterId())
                .startTime(show.getStartTime())
                .price(show.getPrice())
                .build();
    }
}
