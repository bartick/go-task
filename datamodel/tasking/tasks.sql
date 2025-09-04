-- Tasks with parent/child hierarchy
CREATE TABLE tasking.tasks (
  id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  title           VARCHAR(255) NOT NULL,
  description     TEXT,
  status          ENUM('todo','in_progress','done') NOT NULL DEFAULT 'todo',
  priority        TINYINT NOT NULL DEFAULT 0,
  due_date        DATE NULL,
  completed_at    DATETIME NULL,
  parent_task_id  BIGINT UNSIGNED NULL,
  category_id     BIGINT UNSIGNED NULL,
  created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  KEY idx_parent (parent_task_id),
  KEY idx_category (category_id),

  CONSTRAINT fk_task_parent
    FOREIGN KEY (parent_task_id) REFERENCES tasks(id) ON DELETE CASCADE,
  CONSTRAINT fk_task_category
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
) ENGINE=InnoDB;
