import { useMutation, useQuery, useSuspenseQuery } from '@apollo/client';
import * as Dialog from '@radix-ui/react-dialog';
import { Cross2Icon, Pencil2Icon } from '@radix-ui/react-icons';
import * as Label from '@radix-ui/react-label';
import { clsx } from 'clsx';
import {
  ComponentPropsWithoutRef,
  ReactNode,
  useCallback,
  useId,
  useMemo,
  useRef,
  useState,
} from 'react';
import { useDebounce } from 'use-debounce';

import { gql } from '../../../api/__generated__';
import { GetAnkiNoteFieldsAndMappingQuery } from '../../../api/__generated__/graphql';
import { GET_HEALTH_STATUS } from '../../../api/health-status';
import { COLORS } from '../../../colors';
import Button, { ButtonVariant } from '../../../components/Button';
import CodeEditor from '../../../components/CodeEditor';
import { DialogModal, DialogWidth } from '../../../components/DialogModal';
import { LoadingIcon } from '../../../components/Icons/StatusIcon';
import SuspenseLoading from '../../../components/SuspenseLoading';
import Tooltip from '../../../components/Tooltip';
import { useToastify } from '../../../hooks/toastify';
import { apolloErrorToast } from '../../../lib/styled-toast';
import { GET_NOTE_FIELDS_AND_MAPPING } from './api';

export function MappingEdit({ currentNote }: { currentNote?: string }) {
  return (
    <div className="flex flex-col gap-2.5">
      <Label.Root className="text-2xl">Mapping:</Label.Root>
      {currentNote ? (
        <SuspenseLoading>
          <MappingWithFields key={currentNote} />
        </SuspenseLoading>
      ) : (
        <Unavailable />
      )}
    </div>
  );
}

function MappingWithFields() {
  const { data: fieldAndMappingResp } = useSuspenseQuery(GET_NOTE_FIELDS_AND_MAPPING, {
    fetchPolicy: 'network-only',
  });
  if (fieldAndMappingResp.Anki.noteFields.error) {
    return <Unavailable />;
  }
  return <Mapping fieldAndMappingResp={fieldAndMappingResp} />;
}

function Unavailable() {
  return <p className="text-xl">Mapping unavailable</p>;
}

const GET_RENDERED_FIELDS = gql(`
  query RenderFields($fields: [String!]!) {
    RenderFields(fields: $fields) {
      fields {
        result
        error
      }
    }
  }
`);

const SET_ANKI_CONFIG_MAPPING = gql(`
  mutation UpdateMapping($fields: [AnkiConfigMappingElementInput!]!) {
    setAnkiConfigMapping(input: { mapping: $fields }) {
      error {
        fieldErrors {
          key
        }
        valueErrors {
          key
        }
        message
      }
    }
  }
`);

function Mapping({
  fieldAndMappingResp,
}: {
  fieldAndMappingResp: GetAnkiNoteFieldsAndMappingQuery;
}) {
  // we get useful state form response
  const [fields, missingFields] = useMemo(
    () => extractFields(fieldAndMappingResp),
    [fieldAndMappingResp],
  );
  // caching would be nice, but it's not clear how to implement it with apollo
  // easy way to make many requests, but obviously it's not great solution for different reason
  const { data: renderFieldsResponse, loading: renderLoading } = useQuery(GET_RENDERED_FIELDS, {
    variables: {
      fields: fields.map((e) => e.value ?? ''),
    },
  });
  const renderFields = renderFieldsResponse?.RenderFields.fields;
  // editingField is for tracking field that we currently edit
  const [editingField, setEditingField] = useState<MappingField>();
  // updatedFieldName is for tracking lastly updated field (for animation)
  const [updatedFieldName, setUpdatedFieldName] = useState<string>();

  const [setAnkiConfigMapping, { loading: updateInProccess }] = useMutation(
    SET_ANKI_CONFIG_MAPPING,
    {
      refetchQueries: [GET_HEALTH_STATUS, GET_NOTE_FIELDS_AND_MAPPING],
      awaitRefetchQueries: true,
      // we passe empty onError because in this case apollo will not throw in case if network error for some reason
      onError: () => {},
    },
  );

  const toast = useToastify({ type: 'success' });
  const updateFields = useCallback(
    async (fields: MappingField[], missingFields: MappingField[]) => {
      const newFields = fields
        .map((field) => {
          if (!field.value) {
            return null;
          } else {
            return {
              key: field.name,
              value: field.value,
            };
          }
        })
        .filter(<T,>(item: T | null): item is T => item !== null);
      newFields.push(
        ...missingFields.map((field) => ({
          key: field.name,
          value: field.value ?? '',
        })),
      );
      const { data, errors } = await setAnkiConfigMapping({
        variables: {
          fields: newFields,
        },
      });

      if (errors) {
        apolloErrorToast(errors, 'Mapping update failed.', { toast: toast });
        return false;
      } else if (data?.setAnkiConfigMapping.error?.message) {
        toast('Mapping updated failed: ' + data.setAnkiConfigMapping.error.message, {
          type: 'error',
        });
        return false;
      }
      toast('New mapping saved.');
      return true;
    },
    [setAnkiConfigMapping, toast],
  );

  return (
    <>
      <div className="flex flex-col gap-2.5">
        <div className="flex flex-col text-xl gap-8">
          <ul>
            <Field className="pb-2">
              <FieldColumn className="font-bold">Note Field Name</FieldColumn>
              <FieldColumn className="font-bold">Mapping</FieldColumn>
              <FieldColumn className="flex flex-row gap-3 font-bold">
                <div>Rendered Value</div>
                <div className="self-center">
                  {renderLoading && <LoadingIcon className="w-5 h-5" />}
                </div>
              </FieldColumn>
            </Field>
            {fields.map((e, index) => {
              const firstFieldUndefined = index === 0 && !e.value;
              return (
                <Field
                  key={e.name}
                  className={clsx(
                    updatedFieldName == e.name && 'animate-reversePing',
                    // if first field is undefined, require it!
                    firstFieldUndefined && '!bg-light-red',
                  )}
                >
                  <FieldColumn>
                    <div className="flex">
                      <FieldButton
                        tooltip="Edit"
                        aria-label="Edit"
                        disabled={updateInProccess}
                        onClick={() => {
                          setEditingField(e);
                          setUpdatedFieldName(undefined);
                        }}
                      >
                        <Pencil2Icon className="inline" color={COLORS.blue} />
                      </FieldButton>
                      <div className="truncate">{e.name}</div>
                    </div>
                  </FieldColumn>
                  {firstFieldUndefined ? (
                    <FieldColumn className="col-span-2">
                      Error! First field must defined
                    </FieldColumn>
                  ) : (
                    <>
                      <FieldColumn className="font-mono">{e.value}</FieldColumn>
                      <FieldColumn>{renderFields?.at(index)?.result}</FieldColumn>
                    </>
                  )}
                </Field>
              );
            })}
            {missingFields.length > 0 && (
              <Field>
                <FieldColumn className="col-span-3 text-center text-error-red">
                  Following fields present in mapping, but not in Note and therefore must be
                  deleted.
                </FieldColumn>
              </Field>
            )}
            {missingFields.map((e) => {
              return (
                <Field key={e.name}>
                  <FieldColumn className="col-span-3">
                    <div className="flex">
                      <FieldButton
                        tooltip="Delete"
                        aria-label="Delete"
                        disabled={updateInProccess}
                        onClick={async () => {
                          await updateFields(
                            fields,
                            missingFields.filter((field) => field.name != e.name),
                          );
                        }}
                      >
                        <Cross2Icon className="inline" color={COLORS.red} />
                      </FieldButton>
                      <div className="truncate">{e.name}</div>
                    </div>
                  </FieldColumn>
                </Field>
              );
            })}
          </ul>
        </div>
      </div>
      <DialogModal
        widthVariant={DialogWidth.Large}
        open={editingField !== undefined}
        onOpenChange={(open) => open || setEditingField(undefined)}
      >
        {editingField && (
          <EditFieldForm
            key={editingField.name}
            field={editingField}
            updateField={async (field: MappingField, fieldName: string) => {
              const success = await updateFields(
                fields.map((e) => {
                  if (e.name === fieldName) {
                    return field;
                  } else {
                    return e;
                  }
                }),
                missingFields,
              );
              if (success) {
                setUpdatedFieldName(fieldName);
                setEditingField(undefined);
              }
            }}
          />
        )}
      </DialogModal>
    </>
  );
}

function EditFieldForm({
  field,
  updateField,
}: {
  field: MappingField;
  updateField: (field: MappingField, fieldName: string) => Promise<void>;
}) {
  const templateEditorId = useId();
  const [fieldUpdating, setFieldUpdating] = useState(false);
  const update = useCallback(
    async (newTemplate: string) => {
      setFieldUpdating(true);
      await updateField(
        {
          name: field.name,
          value: newTemplate == '' ? undefined : newTemplate,
        },
        field.name,
      );
      setFieldUpdating(false);
    },
    [field, updateField],
  );
  const [template, setTemplate] = useState(field.value ?? '');
  const [debouncedTemplate, debouncedTemplateState] = useDebounce(template, 250);
  const {
    data: currentRenderedFields,
    previousData: previousRenderedFields,
    loading: renderLoading,
  } = useQuery(GET_RENDERED_FIELDS, {
    variables: {
      fields: [debouncedTemplate],
    },
  });
  const renderedField = (currentRenderedFields ?? previousRenderedFields)?.RenderFields.fields[0];
  const renderUpdating = debouncedTemplateState.isPending() || renderLoading;
  const isError = renderUpdating || !!renderedField?.error;
  const updateButtonRef = useRef<HTMLButtonElement>(null);

  return (
    <div className="flex flex-col gap-8">
      <Dialog.Title className="text-2xl font-bold text-blue">
        Field mapping: {field.name}
      </Dialog.Title>

      <div className="flex flex-col gap-2.5">
        <Label.Root className="text-2xl" htmlFor={templateEditorId}>
          Template:
        </Label.Root>
        <CodeEditor
          id={templateEditorId}
          value={template}
          onValueChange={(v) => setTemplate(v)}
          onBlur={() => {
            // for accessibility, when we tab out of editor, we want to focus next button (update)
            updateButtonRef.current?.focus();
          }}
        />
      </div>

      <div className="flex flex-col gap-2.5">
        <Label.Root className="text-2xl">
          <div className="flex gap-2 justify-between place-items-center">
            <div>Result:</div>
            <div>{renderUpdating && <LoadingIcon className="w-5 h-5" />}</div>
          </div>
        </Label.Root>
        <div className="h-40 bg-mid-gray overflow-y-auto whitespace-pre-wrap p-2">
          <div>{renderedField?.result}</div>
          <div className="text-error-red">
            {renderedField?.error && 'Error: ' + renderedField?.error}
          </div>
        </div>
      </div>

      <div className="flex gap-5 max-w-md">
        <Button
          ref={updateButtonRef}
          disabled={fieldUpdating || renderUpdating || isError}
          onClick={() => update(template)}
          className="flex-1"
        >
          Update field
        </Button>
        <Button
          disabled={fieldUpdating}
          className="flex-1"
          variant={ButtonVariant.Dangerous}
          onClick={() => setTemplate(field.value ?? '')}
        >
          Reset
        </Button>
      </div>
    </div>
  );
}

type MappingField = {
  name: string;
  value?: string;
};

type FieldButtonProps = {
  tooltip: string;
} & ComponentPropsWithoutRef<'button'>;

function FieldButton({ children, tooltip, ...rest }: FieldButtonProps) {
  return (
    <Tooltip content={tooltip}>
      <button className="flex items-center justify-center pl-1 pr-2 hover:scale-125" {...rest}>
        {children}
      </button>
    </Tooltip>
  );
}

function Field({ children, className }: { children: ReactNode; className?: string }) {
  return (
    <li
      className={clsx(
        'grid grid-cols-3 grid-cols-[minmax(160px,_15%)_auto_minmax(220px,_30%)] py-1 gap-x-5 even:bg-mid-gray',
        className,
      )}
    >
      {children}
    </li>
  );
}

function FieldColumn({ children, className }: { children: ReactNode; className?: string }) {
  return <div className={clsx('truncate', className)}>{children}</div>;
}

function extractFields(
  resp: GetAnkiNoteFieldsAndMappingQuery,
): [Array<MappingField>, Array<MappingField>] {
  if (!resp.Anki.noteFields.noteFields) {
    return [[], []];
  }
  const fields: Array<MappingField> = resp.Anki.noteFields.noteFields.map((e) => {
    return {
      name: e,
    };
  });
  const missingFields: Array<MappingField> = [];
  resp.AnkiConfig.mapping.forEach((e) => {
    const foundField = fields.find((field) => field.name == e.key);
    if (foundField) {
      foundField.value = e.value;
    } else {
      missingFields.push({
        name: e.key,
        value: e.value,
      });
    }
  });

  return [fields, missingFields];
}
