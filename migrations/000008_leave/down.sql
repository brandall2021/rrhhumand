DELETE FROM permissions WHERE name LIKE 'leave.%';
DROP TABLE IF EXISTS leave_request_history;
DROP TABLE IF EXISTS leave_approvals;
DROP TABLE IF EXISTS leave_requests;
DROP TABLE IF EXISTS leave_balances;
DROP TABLE IF EXISTS holidays;
DROP TABLE IF EXISTS leave_policies;
DROP TABLE IF EXISTS leave_types;
