package com.videoplatform.service;

import com.videoplatform.domain.entity.Danmaku;
import com.videoplatform.domain.entity.User;
import com.videoplatform.domain.entity.Video;
import com.videoplatform.domain.repository.DanmakuRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Service
@RequiredArgsConstructor
public class DanmakuService {

    private final DanmakuRepository danmakuRepository;

    @Transactional
    public Danmaku createDanmaku(Long userId, Long videoId, String content, Double positionSeconds, 
                                  Danmaku.DanmakuStyle style, String color, Integer fontSize) {
        Danmaku danmaku = Danmaku.builder()
            .user(createUserRef(userId))
            .video(createVideoRef(videoId))
            .content(content)
            .positionSeconds(positionSeconds)
            .style(style != null ? style : Danmaku.DanmakuStyle.SCROLL)
            .color(color != null ? color : "#FFFFFF")
            .fontSize(fontSize != null ? fontSize : 24)
            .build();
        
        return danmakuRepository.save(danmaku);
    }

    public List<Danmaku> getDanmakusByVideo(Long videoId, Double startTime, Double endTime) {
        if (startTime != null && endTime != null) {
            return danmakuRepository.findByVideoIdAndPositionRange(videoId, startTime, endTime);
        }
        return danmakuRepository.findByVideoIdOrderByPositionSeconds(videoId);
    }

    public Long getDanmakuCount(Long videoId) {
        return danmakuRepository.countByVideoId(videoId);
    }

    @Transactional
    public void deleteDanmaku(Long danmakuId) {
        danmakuRepository.deleteById(danmakuId);
    }

    private User createUserRef(Long id) {
        User user = new User();
        user.setId(id);
        return user;
    }

    private Video createVideoRef(Long id) {
        Video video = new Video();
        video.setId(id);
        return video;
    }
}
