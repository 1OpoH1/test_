-- ============================================================
-- TABLE: users
-- Stores registered users.
-- ============================================================

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR UNIQUE NOT NULL,
    username VARCHAR NOT NULL,
    completed_chats INTEGER DEFAULT 0
);


-- ============================================================
-- TABLE: chats
-- Stores available chat scenarios.
-- ============================================================

CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    description TEXT,
    difficulty VARCHAR NOT NULL,
    role VARCHAR NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


-- ============================================================
-- TABLE: achievements
-- Stores available achievements.
-- ============================================================

CREATE TABLE achievements (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    description TEXT,
    icon VARCHAR,
    condition_type VARCHAR NOT NULL,
    condition_value VARCHAR NOT NULL
);


-- ============================================================
-- TABLE: statistics
-- Stores aggregated user statistics.
-- One user has one statistics record.
-- ============================================================

CREATE TABLE statistics (
    id SERIAL PRIMARY KEY,

    user_id INTEGER NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    total_chats INTEGER DEFAULT 0,
    completed_chats INTEGER DEFAULT 0,
    failed_chats INTEGER DEFAULT 0,
    total_messages INTEGER DEFAULT 0,

    time_spent INTEGER DEFAULT 0,
    -- Stored in seconds

    success_rate NUMERIC(5, 2) DEFAULT 0.00,
    -- Percentage rate

    best_score INTEGER DEFAULT 0,

    UNIQUE (user_id)
);


-- ============================================================
-- TABLE: user_achievements
-- Junction table between users and achievements.
-- ============================================================

CREATE TABLE user_achievements (
    id SERIAL PRIMARY KEY,

    user_id INTEGER NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    achievement_id INTEGER NOT NULL
        REFERENCES achievements(id)
        ON DELETE CASCADE,

    received_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (user_id, achievement_id)
);


-- ============================================================
-- TABLE: chat_steps
-- Stores ordered messages that make up a chat scenario.
-- ============================================================

CREATE TABLE chat_steps (
    id SERIAL PRIMARY KEY,

    chat_id INTEGER NOT NULL
        REFERENCES chats(id)
        ON DELETE CASCADE,

    step_number INTEGER NOT NULL,

    role VARCHAR NOT NULL,

    message_text TEXT NOT NULL,

    response_type VARCHAR
    -- e.g. 'OPTION', 'FREE_TEXT', 'END'
);


-- Step number must be unique inside one chat.
CREATE UNIQUE INDEX idx_chat_step_number
    ON chat_steps (chat_id, step_number);


-- ============================================================
-- TABLE: chat_options
-- Stores answer options and scoring for specific steps.
-- ============================================================

CREATE TABLE chat_options (
    id SERIAL PRIMARY KEY,

    chat_id INTEGER NOT NULL
        REFERENCES chats(id)
        ON DELETE CASCADE,

    step_number INTEGER NOT NULL,

    option_text TEXT NOT NULL,

    is_correct BOOLEAN NOT NULL,

    explanation TEXT,

    points INTEGER DEFAULT 0
);


-- ============================================================
-- TABLE: chat_sessions
-- Represents a user's attempt to complete a chat.
-- ============================================================

CREATE TABLE chat_sessions (
    id SERIAL PRIMARY KEY,

    user_id INTEGER NOT NULL
        REFERENCES users(id)
        ON DELETE RESTRICT,

    chat_id INTEGER NOT NULL
        REFERENCES chats(id)
        ON DELETE RESTRICT,

    status VARCHAR NOT NULL,
    -- e.g. 'IN_PROGRESS', 'COMPLETED', 'ABANDONED'

    started_at TIMESTAMP WITH TIME ZONE
        DEFAULT CURRENT_TIMESTAMP,

    finished_at TIMESTAMP WITH TIME ZONE,

    score INTEGER DEFAULT 0,

    current_step INTEGER DEFAULT 0
);


-- ============================================================
-- TABLE: messages
-- Stores all messages exchanged during a chat session.
-- ============================================================

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,

    session_id INTEGER NOT NULL
        REFERENCES chat_sessions(id)
        ON DELETE CASCADE,

    role VARCHAR NOT NULL,
    -- e.g. 'USER', 'NPC', 'SYSTEM'

    message TEXT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE
        DEFAULT CURRENT_TIMESTAMP
);


-- ============================================================
-- TABLE: chat_keywords
-- Stores keywords associated with chat steps.
-- Used for evaluation of free-text answers.
-- ============================================================

CREATE TABLE chat_keywords (
    id SERIAL PRIMARY KEY,

    chat_id INTEGER NOT NULL
        REFERENCES chats(id)
        ON DELETE CASCADE,

    step_number INTEGER NOT NULL,

    keyword VARCHAR(255) NOT NULL,

    weight FLOAT DEFAULT 1.0,

    is_required BOOLEAN DEFAULT FALSE
);


-- ============================================================
-- TABLE: session_answers
-- Stores answers given by the user during a chat session.
-- ============================================================

CREATE TABLE session_answers (
    id SERIAL PRIMARY KEY,

    session_id INTEGER NOT NULL
        REFERENCES chat_sessions(id)
        ON DELETE CASCADE,

    step_id INTEGER NOT NULL
        REFERENCES chat_steps(id)
        ON DELETE CASCADE,

    option_id INTEGER
        REFERENCES chat_options(id)
        ON DELETE SET NULL,

    free_text TEXT,

    is_correct BOOLEAN NOT NULL DEFAULT FALSE,

    ai_evaluation TEXT,

    created_at TIMESTAMP WITH TIME ZONE
        DEFAULT CURRENT_TIMESTAMP
);


-- ============================================================
-- TABLE: lessons
-- Stores educational lessons.
-- ============================================================

CREATE TABLE lessons (
    id SERIAL PRIMARY KEY,

    title VARCHAR(255) NOT NULL,

    description TEXT,

    theory_content TEXT NOT NULL,

    icon VARCHAR(255),

    "order" INTEGER DEFAULT 0,

    is_active BOOLEAN DEFAULT TRUE,

    passing_score INTEGER NOT NULL
);


-- ============================================================
-- TABLE: lesson_questions
-- Stores questions belonging to lessons.
-- ============================================================

CREATE TABLE lesson_questions (
    id SERIAL PRIMARY KEY,

    lesson_id INTEGER NOT NULL
        REFERENCES lessons(id)
        ON DELETE CASCADE,

    question_text TEXT NOT NULL,

    explanation TEXT,

    "order" INTEGER NOT NULL
);


-- Each question order must be unique inside a lesson.
CREATE UNIQUE INDEX idx_lesson_question_order
    ON lesson_questions (lesson_id, "order");


-- ============================================================
-- TABLE: lesson_options
-- Stores answer options for lesson questions.
-- ============================================================

CREATE TABLE lesson_options (
    id SERIAL PRIMARY KEY,

    question_id INTEGER NOT NULL
        REFERENCES lesson_questions(id)
        ON DELETE CASCADE,

    option_text TEXT NOT NULL,

    is_correct BOOLEAN NOT NULL,

    "order" INTEGER NOT NULL
);


-- Each option order must be unique inside one question.
CREATE UNIQUE INDEX idx_lesson_option_order
    ON lesson_options (question_id, "order");


-- ============================================================
-- TABLE: user_lessons
-- Stores user's progress through lessons.
-- ============================================================

CREATE TABLE user_lessons (
    id SERIAL PRIMARY KEY,

    user_id INTEGER NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    lesson_id INTEGER NOT NULL
        REFERENCES lessons(id)
        ON DELETE CASCADE,

    theory_read BOOLEAN DEFAULT FALSE,

    quiz_score INTEGER DEFAULT 0,

    quiz_passed BOOLEAN DEFAULT FALSE,

    chats_completed INTEGER DEFAULT 0,

    started_at TIMESTAMP WITH TIME ZONE,

    completed_at TIMESTAMP WITH TIME ZONE,

    UNIQUE (user_id, lesson_id)
);


-- ============================================================
-- TABLE: lesson_chats
-- Links lessons with chats.
-- Many-to-many relationship between lessons and chats.
-- ============================================================

CREATE TABLE lesson_chats (
    lesson_id INTEGER NOT NULL
        REFERENCES lessons(id)
        ON DELETE CASCADE,

    chat_id INTEGER NOT NULL
        REFERENCES chats(id)
        ON DELETE CASCADE,

    "order" INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY (lesson_id, chat_id)
);


-- ============================================================
-- INDEXES
-- ============================================================

CREATE INDEX idx_statistics_user_id
    ON statistics(user_id);

CREATE INDEX idx_user_achievements_user_id
    ON user_achievements(user_id);

CREATE INDEX idx_user_achievements_achievement_id
    ON user_achievements(achievement_id);

CREATE INDEX idx_chat_steps_chat_id
    ON chat_steps(chat_id);

CREATE INDEX idx_chat_options_chat_id
    ON chat_options(chat_id);

CREATE INDEX idx_chat_sessions_user_id
    ON chat_sessions(user_id);

CREATE INDEX idx_chat_sessions_chat_id
    ON chat_sessions(chat_id);

CREATE INDEX idx_messages_session_id
    ON messages(session_id);

CREATE INDEX idx_chat_keywords_chat_id
    ON chat_keywords(chat_id);

CREATE INDEX idx_session_answers_session_id
    ON session_answers(session_id);

CREATE INDEX idx_session_answers_step_id
    ON session_answers(step_id);

CREATE INDEX idx_session_answers_option_id
    ON session_answers(option_id);

CREATE INDEX idx_lesson_questions_lesson_id
    ON lesson_questions(lesson_id);

CREATE INDEX idx_lesson_options_question_id
    ON lesson_options(question_id);

CREATE INDEX idx_user_lessons_user_id
    ON user_lessons(user_id);

CREATE INDEX idx_user_lessons_lesson_id
    ON user_lessons(lesson_id);

CREATE INDEX idx_lesson_chats_chat_id
    ON lesson_chats(chat_id);