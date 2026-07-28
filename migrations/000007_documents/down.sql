DROP TABLE IF EXISTS document_shares;
DROP TABLE IF EXISTS document_access_logs;
DROP TABLE IF EXISTS document_tag_relations;
DROP TABLE IF EXISTS document_tags;
DROP TABLE IF EXISTS document_permissions;
DROP TABLE IF EXISTS document_versions;
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS document_categories;

DELETE FROM permissions WHERE resource = 'documents';
