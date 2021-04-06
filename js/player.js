function isInChineseMainland(callback) {
    var url = 'https://graph.facebook.com/feed?callback=t';
    var xhr = new XMLHttpRequest();
    var called = false;
    xhr.open('GET', url);
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4 && xhr.status === 200) {
            called = true;
            callback(false);
        }
    };
    xhr.send();
    setTimeout(function() {
        if (!called) {
            xhr.abort();
            callback(true);
        }
    }, 1000);
}

function setPlayer(in_chinese) {
    document.querySelectorAll('.player').forEach(function(frame) {
        if (in_chinese) {
            var video_id = frame.getAttribute('bilibili');
            frame.setAttribute('src', '//player.bilibili.com/player.html?bvid=' + video_id);
        } else {
            var video_id = frame.getAttribute('youtube');
            frame.setAttribute('src', '//www.youtube.com/embed/' + video_id);
        }
        frame.removeAttribute('youtube');
        frame.removeAttribute('bilibili');
        frame.setAttribute('width', '880');
        frame.setAttribute('height', '495');
        frame.setAttribute('frameborder', '0');
        frame.setAttribute('allow',
            'accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture');
    });
    if (in_chinese) {
        ;
    } else {
    }
}

document.addEventListener('DOMContentLoaded', function(event) {
    if (document.querySelectorAll('.player').length > 0) {
        isInChineseMainland(setPlayer);
    }
});
