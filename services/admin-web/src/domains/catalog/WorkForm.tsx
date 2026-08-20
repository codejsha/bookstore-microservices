import { Button } from "@bookstore/design/ui/button";
import { Input } from "@bookstore/design/ui/input";
import { Textarea } from "@bookstore/design/ui/textarea";
import { useForm } from "@tanstack/react-form";
import { FormActions, FormField } from "@/shared/components/ui/form";
import type { Author, Work } from "./types";

export interface WorkFormValues {
  title: string;
  description: string;
  first_publish_date: string;
  ol_key: string;
  author_uids: string[];
  subject_names: string[];
}

interface WorkFormProps {
  authors: Author[];
  work?: Work;
  submitLabel: string;
  pending: boolean;
  onSubmit: (values: WorkFormValues) => void;
  onCancel: () => void;
}

function initialValues(work: Work | undefined): WorkFormValues {
  return {
    title: work?.title ?? "",
    description: work?.description ?? "",
    first_publish_date: work?.first_publish_date ?? "",
    ol_key: work?.ol_key ?? "",
    author_uids: work?.authors.map((a) => a.uid) ?? [],
    subject_names: work?.subjects.map((s) => s.name) ?? [],
  };
}

export function WorkForm({
  authors,
  work,
  submitLabel,
  pending,
  onSubmit,
  onCancel,
}: WorkFormProps) {
  const form = useForm({
    defaultValues: initialValues(work),
    onSubmit: ({ value }) => onSubmit(value),
  });

  return (
    <form
      className="flex max-w-2xl flex-col gap-4"
      onSubmit={(e) => {
        e.preventDefault();
        form.handleSubmit();
      }}
    >
      <form.Field
        name="title"
        validators={{
          onSubmit: ({ value }) =>
            value.trim() ? undefined : "Title is required",
        }}
      >
        {(field) => (
          <FormField label="Title" errors={field.state.meta.errors}>
            {({ id, describedBy }) => (
              <Input
                id={id}
                aria-describedby={describedBy}
                value={field.state.value}
                onBlur={field.handleBlur}
                onChange={(e) => field.handleChange(e.target.value)}
              />
            )}
          </FormField>
        )}
      </form.Field>

      <form.Field
        name="author_uids"
        validators={{
          onSubmit: ({ value }) =>
            value.length > 0 ? undefined : "Pick at least one author",
        }}
      >
        {(field) => (
          <FormField
            label="Authors"
            errors={field.state.meta.errors}
            hint="Authors must already exist in the catalog."
          >
            {({ id, describedBy }) => (
              <select
                id={id}
                aria-describedby={describedBy}
                multiple
                size={6}
                className="rounded-md border border-input bg-background p-2 text-sm"
                value={field.state.value}
                onBlur={field.handleBlur}
                onChange={(e) =>
                  field.handleChange(
                    Array.from(e.target.selectedOptions, (o) => o.value),
                  )
                }
              >
                {authors.map((author) => (
                  <option key={author.uid} value={author.uid}>
                    {author.name}
                  </option>
                ))}
              </select>
            )}
          </FormField>
        )}
      </form.Field>

      <form.Field name="subject_names">
        {(field) => (
          <FormField
            label="Subjects"
            errors={field.state.meta.errors}
            hint="Comma separated. Unlike authors, subjects are created on demand from their names."
          >
            {({ id, describedBy }) => (
              <Input
                id={id}
                aria-describedby={describedBy}
                value={field.state.value.join(", ")}
                onBlur={field.handleBlur}
                onChange={(e) =>
                  field.handleChange(
                    e.target.value
                      .split(",")
                      .map((s) => s.trim())
                      .filter(Boolean),
                  )
                }
              />
            )}
          </FormField>
        )}
      </form.Field>

      <form.Field name="first_publish_date">
        {(field) => (
          <FormField
            label="First published"
            errors={field.state.meta.errors}
            hint="Free text, as the catalog stores it: 1965, 1965-06, or a full date."
          >
            {({ id, describedBy }) => (
              <Input
                id={id}
                aria-describedby={describedBy}
                value={field.state.value}
                onBlur={field.handleBlur}
                onChange={(e) => field.handleChange(e.target.value)}
              />
            )}
          </FormField>
        )}
      </form.Field>

      <form.Field name="ol_key">
        {(field) => (
          <FormField label="Open Library key" errors={field.state.meta.errors}>
            {({ id, describedBy }) => (
              <Input
                id={id}
                aria-describedby={describedBy}
                value={field.state.value}
                onBlur={field.handleBlur}
                onChange={(e) => field.handleChange(e.target.value)}
              />
            )}
          </FormField>
        )}
      </form.Field>

      <form.Field name="description">
        {(field) => (
          <FormField label="Description" errors={field.state.meta.errors}>
            {({ id, describedBy }) => (
              <Textarea
                id={id}
                aria-describedby={describedBy}
                rows={5}
                value={field.state.value}
                onBlur={field.handleBlur}
                onChange={(e) => field.handleChange(e.target.value)}
              />
            )}
          </FormField>
        )}
      </form.Field>

      <FormActions>
        <Button
          type="button"
          variant="outline"
          onClick={onCancel}
          disabled={pending}
        >
          Cancel
        </Button>
        <Button type="submit" disabled={pending}>
          {pending ? "Saving…" : submitLabel}
        </Button>
      </FormActions>
    </form>
  );
}
