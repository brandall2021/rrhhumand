-- 000026: Down - remove missing worker tables

DROP TABLE IF EXISTS training_events CASCADE;
DROP TABLE IF EXISTS training_course_offerings CASCADE;
DROP TABLE IF EXISTS scoring_matches CASCADE;
