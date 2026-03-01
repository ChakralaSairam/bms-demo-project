package com.test.user.service;

import com.test.user.dto.MovieDTO;
import com.test.user.models.Movie;
import com.test.user.repository.MovieRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Service
@RequiredArgsConstructor
@Slf4j
public class MovieAdminService {

    private final MovieRepository movieRepository;

    @Transactional
    public MovieDTO createMovie(MovieDTO dto) {
        log.info("Creating movie: {}", dto.getTitle());
        Movie movie = Movie.builder()
                .title(dto.getTitle())
                .genre(dto.getGenre())
                .duration(dto.getDuration())
                .language(dto.getLanguage())
                .build();
        movie = movieRepository.save(movie);
        return toDTO(movie);
    }

    @Transactional
    public MovieDTO updateMovie(Long id, MovieDTO dto) {
        log.info("Updating movie ID: {}", id);
        Movie movie = movieRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Movie not found"));
        movie.setTitle(dto.getTitle());
        movie.setGenre(dto.getGenre());
        movie.setDuration(dto.getDuration());
        movie.setLanguage(dto.getLanguage());
        movie = movieRepository.save(movie);
        return toDTO(movie);
    }

    @Transactional
    public void deleteMovie(Long id) {
        log.info("Deleting movie ID: {}", id);
        movieRepository.deleteById(id);
    }

    private MovieDTO toDTO(Movie movie) {
        return MovieDTO.builder()
                .id(movie.getId())
                .title(movie.getTitle())
                .genre(movie.getGenre())
                .duration(movie.getDuration())
                .language(movie.getLanguage())
                .build();
    }
}
