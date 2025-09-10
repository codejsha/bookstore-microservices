import type { WorkCreateRequest, WorkUpdateRequest } from "./types";
import type { WorkFormValues } from "./WorkForm";

function blankToUndefined(value: string): string | undefined {
  const trimmed = value.trim();
  return trimmed === "" ? undefined : trimmed;
}

export function toCreateRequest(values: WorkFormValues): WorkCreateRequest {
  return {
    title: values.title.trim(),
    description: blankToUndefined(values.description),
    first_publish_date: blankToUndefined(values.first_publish_date),
    ol_key: blankToUndefined(values.ol_key),
    author_uids: values.author_uids,
    subject_names:
      values.subject_names.length > 0 ? values.subject_names : undefined,
  };
}

export function toUpdateRequest(values: WorkFormValues): WorkUpdateRequest {
  return {
    title: blankToUndefined(values.title),
    description: blankToUndefined(values.description),
    first_publish_date: blankToUndefined(values.first_publish_date),
    ol_key: blankToUndefined(values.ol_key),
    author_uids: values.author_uids,
    subject_names: values.subject_names,
  };
}
