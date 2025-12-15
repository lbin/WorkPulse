import { client } from "./client";

export interface MeetingPayload {
  title: string;
  agenda?: string;
  scheduled_at: string;
  duration_minutes?: number;
  facilitator_user_id?: string | null;
  attendee_ids?: string[];
  status?: string;
  team_id?: string | null;
  notes?: string;
}

export function listMeetings(params: { start?: string; end?: string; team_id?: string }) {
  return client.get("/meetings", { params });
}

export function getMeeting(id: string) {
  return client.get(`/meetings/${id}`);
}

export function createMeeting(payload: MeetingPayload) {
  return client.post("/meetings", payload);
}

export function updateMeeting(id: string, payload: Partial<MeetingPayload>) {
  return client.put(`/meetings/${id}`, payload);
}

export function addMeetingAction(meetingId: string, payload: { title: string; owner_user_id?: string; due_date?: string; status?: string; related_task_id?: string }) {
  return client.post(`/meetings/${meetingId}/actions`, payload);
}

export function assignMeetingAction(actionId: string, payload: { owner_user_id?: string; status?: string; related_task_id?: string }) {
  return client.post(`/meetings/actions/${actionId}/assign`, payload);
}

export function convertMeetingAction(actionId: string) {
  return client.post(`/meetings/actions/${actionId}/convert-task`);
}

export function addMeetingLink(meetingId: string, payload: { target_type: string; target_id: string; relation?: string }) {
  return client.post(`/meetings/${meetingId}/links`, payload);
}
