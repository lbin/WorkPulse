import type { AxiosResponse } from "axios";
import { api } from "./client";

export interface Meeting {
  id: string;
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

export type MeetingPayload = Omit<Meeting, "id">;

export interface MeetingAction {
  id: string;
  title: string;
  owner_user_id?: string;
  due_date?: string;
  status?: string;
}

export interface MeetingLink {
  id: string;
  target_type: string;
  target_id: string;
  relation?: string;
}

type MeetingListResponse = AxiosResponse<{ data: Meeting[] }>;

type MeetingDetailResponse = AxiosResponse<{
  data: { meeting: Meeting; actions: MeetingAction[]; links: MeetingLink[] };
}>;

export function listMeetings(params: { start?: string; end?: string; team_id?: string }) {
  return api.get<MeetingListResponse["data"]>("/meetings", { params });
}

export function getMeeting(id: string) {
  return api.get<MeetingDetailResponse["data"]>(`/meetings/${id}`);
}

export function createMeeting(payload: MeetingPayload) {
  return api.post("/meetings", payload);
}

export function updateMeeting(id: string, payload: Partial<MeetingPayload>) {
  return api.put(`/meetings/${id}`, payload);
}

export function addMeetingAction(meetingId: string, payload: { title: string; owner_user_id?: string; due_date?: string; status?: string; related_task_id?: string }) {
  return api.post(`/meetings/${meetingId}/actions`, payload);
}

export function assignMeetingAction(actionId: string, payload: { owner_user_id?: string; status?: string; related_task_id?: string }) {
  return api.post(`/meetings/actions/${actionId}/assign`, payload);
}

export function convertMeetingAction(actionId: string) {
  return api.post(`/meetings/actions/${actionId}/convert-task`);
}

export function addMeetingLink(meetingId: string, payload: { target_type: string; target_id: string; relation?: string }) {
  return api.post(`/meetings/${meetingId}/links`, payload);
}
