DROP TABLE IF EXISTS survey_answer_options;
DROP TABLE IF EXISTS survey_answers;
DROP TABLE IF EXISTS survey_responses;
DROP TABLE IF EXISTS survey_targets;
DROP TABLE IF EXISTS survey_options;
DROP TABLE IF EXISTS survey_questions;
DROP TABLE IF EXISTS surveys;

DELETE FROM permissions WHERE resource = 'surveys';
