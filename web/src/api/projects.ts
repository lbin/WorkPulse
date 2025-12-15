import { api } from "./client";

export interface Project {
  id: string;
  name: string;
  description: string;
  status: string;
  start_date?: string | null;
  end_date?: string | null;
  milestone_count: number;
  task_count: number;
  progress: number;
}

export interface Milestone {
  id: string;
  title: string;
  description: string;
  due_date?: string | null;
  status: string;
  progress: number;
}

export interface Task {
  id: string;
  title: string;
  status: string;
  milestone_id?: string | null;
  parent_task_id?: string | null;
  tags: string[];
}

export interface TaskLink {
  id: string;
  entity_type: string;
  entity_id: string;
}

export async function fetchProjects() {
  const { data } = await api.get<{ data: Project[] }>("/projects");
  return data.data;
}

export async function fetchMilestones(projectId: string) {
  const { data } = await api.get<{ data: Milestone[] }>(`/projects/${projectId}/milestones`);
  return data.data;
}

export async function fetchTasks(projectId: string) {
  const { data } = await api.get<{ data: Task[] }>(`/projects/${projectId}/tasks`);
  return data.data;
}

export async function fetchTaskLinks(taskId: string) {
  const { data } = await api.get<{ data: TaskLink[] }>(`/tasks/${taskId}/links`);
  return data.data;
}

export async function linkTask(taskId: string, entityType: string, entityId: string) {
  const { data } = await api.post<{ data: TaskLink }>(`/tasks/${taskId}/links`, {
    entity_type: entityType,
    entity_id: entityId,
  });
  return data.data;
}

export async function unlinkTask(taskId: string, linkId: string) {
  await api.delete(`/tasks/${taskId}/links/${linkId}`);
}

export async function fetchProjectOKRProgress(projectId: string) {
  const { data } = await api.get<{ data: { project_progress: number; key_results: { key_result_id: string; progress: number }[] } }>(
    `/projects/${projectId}/okr-progress`
  );
  return data.data;
}
