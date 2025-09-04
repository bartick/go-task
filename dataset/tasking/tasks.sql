-- Root tasks
INSERT INTO tasking.tasks (title, status, priority, due_date, category_id)
VALUES
  ('Implémenter API Auth', 'in_progress', 3, '2025-08-20', 1),
  ('Créer Dashboard UI', 'todo', 2, '2025-08-25', 2),
  ('Corriger bug login', 'todo', 4, '2025-08-15', 3);

-- Subtasks for task 1
INSERT INTO tasks (title, status, priority, due_date, parent_task_id, category_id)
VALUES
  ('Ajouter OAuth2', 'todo', 2, '2025-08-18', 1, 1),
  ('Configurer JWT', 'todo', 3, '2025-08-19', 1, 1);

-- Subtasks for task 2
INSERT INTO tasks (title, status, priority, due_date, parent_task_id, category_id)
VALUES
  ('Créer composant Graphique', 'in_progress', 2, '2025-08-22', 2, 2);