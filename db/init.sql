CREATE TABLE IF NOT EXISTS User (
    usr_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) NOT NULL UNIQUE,
    password VARCHAR(72) NOT NULL, -- bcryptの最大長が72らしい
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Tweet (
    tw_id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    usr_id INT NOT NULL,
    body VARCHAR(280) NOT NULL,
    replied_to BIGINT UNSIGNED NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_replied_to(replied_to),
);

CREATE TABLE IF NOT EXISTS Favorite (
    tw_id BIGINT UNSIGNED NOT NULL,
    usr_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tw_id, usr_id)
);

CREATE TABLE IF NOT EXISTS TweetStats (
    tw_id BIGINT UNSIGNED PRIMARY KEY,
    fav_count INT NOT NULL DEFAULT 0,
    reply_count INT NOT NULL DEFAULT 0
);

INSERT INTO TweetStats (tw_id) SELECT tw_id FROM Tweet; 
UPDATE TweetStats T JOIN (SELECT tw_id, count(1) fav_count FROM Favorite GROUP BY tw_id) F USING(tw_id) SET T.fav_count=F.fav_count;
UPDATE TweetStats T JOIN (SELECT replied_to tw_id, count(1) reply_count FROM Tweet GROUP BY replied_to) F USING(tw_id) SET T.reply_count=F.reply_count;

CREATE TRIGGER trigger_create_stats
    AFTER INSERT
    ON Tweet FOR EACH ROW 
        INSERT INTO TweetStats (tw_id) VALUES (NEW.tw_id);

CREATE TRIGGER trigger_addcount_reply
    AFTER INSERT
    ON Tweet FOR EACH ROW 
        UPDATE TweetStats SET reply_count=reply_count+1 WHERE tw_id=NEW.replied_to;

CREATE TRIGGER trigger_addcount_fav
    AFTER INSERT
    ON Favorite FOR EACH ROW
    UPDATE TweetStats SET fav_count=fav_count+1 WHERE tw_id=NEW.tw_id;

CREATE TRIGGER trigger_subcount_fav
    AFTER DELETE
    ON Favorite FOR EACH ROW
    UPDATE TweetStats SET fav_count=fav_count-1 WHERE tw_id=OLD.tw_id;
