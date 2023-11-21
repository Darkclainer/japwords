import { useMutation, useQuery, useSuspenseQuery } from '@apollo/client';
import { Cross2Icon, Pencil2Icon } from '@radix-ui/react-icons';
import * as Dialog from '@radix-ui/react-dialog';
import * as Label from '@radix-ui/react-label';
import {
  ComponentPropsWithoutRef,
  Dispatch,
  ReactNode,
  SetStateAction,
  useCallback,
  useId,
  useMemo,
  useState,
} from 'react';

import { gql } from '../../../api/__generated__';
import { GetAnkiNoteFieldsAndMappingQuery } from '../../../api/__generated__/graphql';
import { COLORS } from '../../../colors';
import Button, { ButtonVariant } from '../../../components/Button';
import SuspenseLoading from '../../../components/SuspenseLoading';
import { clsx } from 'clsx';
import { LoadingIcon } from '../../../components/StatusIcon';
import Tooltip from '../../../components/Tooltip';
import { DialogModal, DialogWidth } from '../../../components/DialogModal';
import CodeEditor from '../../../components/CodeEditor';
import { useDebounce } from 'use-debounce';
import { GET_HEALTH_STATUS } from '../../../api/health-status';
import { useToastify } from '../../../hooks/toastify';
import { apolloErrorToast } from '../../../lib/styled-toast';

export function MappingEdit({ currentNote }: { currentNote?: string }) {
  return (
    <div className="flex flex-col gap-2.5">
      <Label.Root className="text-2xl">Mapping:</Label.Root>
      {currentNote ? (
        <SuspenseLoading>
          <MappingWithFields key={currentNote} currentNote={currentNote} />
        </SuspenseLoading>
      ) : (
        <Unavailable />
      )}
    </div>
  );
}

function MappingWithFields({ currentNote }: { currentNote: string }) {
  const { data: fieldAndMappingResp } = useSuspenseQuery(GET_NOTE_FIELDS_AND_MAPPING, {
    fetchPolicy: 'no-cache',
    variables: {
      noteName: currentNote,
    },
  });
  if (fieldAndMappingResp.Anki.error) {
    return <Unavailable />;
  }
  return <Mapping fieldAndMappingResp={fieldAndMappingResp} />;
}

function Unavailable() {
  return <p className="text-xl">Mapping unavailable</p>;
}

const GET_NOTE_FIELDS_AND_MAPPING = gql(`
  query GetAnkiNoteFieldsAndMapping($noteName: String!) {
    AnkiConfig {
      mapping {
        key
        value
      }
    }
    Anki {
      anki {
        noteFields(name: $noteName)
      }
      error {
        ... on Error {
          message
        }
      }
    }
  }
`);

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
  const [initalFields, initialMissingFields] = useMemo(
    () => extractFields(fieldAndMappingResp),
    [fieldAndMappingResp],
  );
  // get prerendered values for inital fields only
  const { data: renderFieldsResponse, loading: renderLoading } = useQuery(GET_RENDERED_FIELDS, {
    variables: {
      fields: initalFields.map((e) => e.value ?? ''),
    },
  });
  const renderFields = renderFieldsResponse?.RenderFields.fields;
  // get seperate state for fields and missing fields so that we can alter it and save later
  const [fields, setFields] = useState(initalFields);
  const [missingFields, setMissingFields] = useState(initialMissingFields);
  // editingFieldId is for tracking field that we currently edit
  const [editingFieldId, setEditingFieldId] = useState<number>();
  // updatedFieldId is for tracking lastly updated field (for animation)
  const [updatedFieldId, setUpdatedFieldId] = useState<number>();

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
  const updateFields = useCallback(async (fields: MappingField[]) => {
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
    const { data, errors } = await setAnkiConfigMapping({
      variables: {
        fields: newFields,
      },
    });

    if (errors) {
      apolloErrorToast(errors, 'Mapping update failed.', { toast: toast });
      return;
    } else if (data?.setAnkiConfigMapping.error?.message) {
      toast('Mapping updated failed: ' + data.setAnkiConfigMapping.error.message, {
        type: 'error',
      });
      return;
    }
    toast('New mapping saved.');
  }, []);

  const canBeUpdated = useMemo(() => {
    // safe check, actually should not be important
    if (initalFields.length != fields.length) {
      return false;
    }
    if (missingFields.length !== initialMissingFields.length) {
      return true;
    }
    return !initalFields.every((field, i) => field.value === fields[i].value);
  }, [initalFields, fields, initialMissingFields, missingFields]);
  const hasError = missingFields.length !== 0;

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
                <div className="self-center">{renderLoading && <LoadingIcon size="1.25rem" />}</div>
              </FieldColumn>
            </Field>
            {fields.map((e, index) => (
              <Field
                key={e.name}
                className={clsx(updatedFieldId == index && 'animate-reversePing')}
              >
                <FieldColumn>
                  <div className="flex">
                    <FieldButton
                      tooltip="Edit"
                      aria-label="Edit"
                      onClick={() => {
                        setEditingFieldId(index);
                        setUpdatedFieldId(undefined);
                      }}
                    >
                      <Pencil2Icon className="inline" color={COLORS.blue} />
                    </FieldButton>
                    <div className="truncate">{e.name}</div>
                  </div>
                </FieldColumn>
                <FieldColumn className="font-mono">{e.value}</FieldColumn>
                <FieldColumn>
                  {e.example === undefined ? renderFields?.at(index)?.result : e.example}
                </FieldColumn>
              </Field>
            ))}
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
                        onClick={() => {
                          setMissingFields((fields) => fields?.filter((value) => value != e));
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
        <div className="flex flex-row max-w-md gap-5">
          <Button
            className="flex-1"
            disabled={hasError || updateInProccess || !canBeUpdated}
            onClick={() => updateFields(fields)}
          >
            Update mapping
          </Button>
          <Button
            className="flex-1"
            variant={ButtonVariant.Dangerous}
            disabled={updateInProccess}
            onClick={() => {
              setFields(initalFields);
              setMissingFields(initialMissingFields);
              setEditingFieldId(undefined);
              setUpdatedFieldId(undefined);
            }}
          >
            Reset
          </Button>
        </div>
      </div>
      <DialogModal
        widthVariant={DialogWidth.Large}
        open={editingFieldId !== undefined}
        onOpenChange={(open) => open || setEditingFieldId(undefined)}
      >
        <EditFieldForm
          key={editingFieldId}
          fieldId={editingFieldId}
          fields={fields}
          setFields={(...args) => {
            setFields(...args);
            setUpdatedFieldId(editingFieldId);
            setEditingFieldId(undefined);
          }}
        />
      </DialogModal>
    </>
  );
}

function EditFieldForm({
  fieldId,
  fields,
  setFields,
}: {
  fieldId?: number;
  fields: Array<MappingField>;
  setFields: Dispatch<SetStateAction<Array<MappingField>>>;
}) {
  const templateEditorId = useId();
  const field = fieldId !== undefined ? fields[fieldId] : { name: '' };
  const updateField = useCallback(
    (newTemplate: string, example?: string) =>
      setFields((fields) =>
        fields.map((item, i) => {
          if (i === fieldId) {
            return {
              name: field.name,
              value: newTemplate == '' ? undefined : newTemplate,
              example: example,
            };
          } else {
            return item;
          }
        }),
      ),
    [fieldId, fields, setFields],
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
  const updating = debouncedTemplateState.isPending() || renderLoading;
  const isError = updating || !!renderedField?.error;
  return (
    <div className="flex flex-col gap-8">
      <Dialog.Title className="text-2xl font-bold text-blue">
        Field mapping: {field.name}
      </Dialog.Title>

      <div className="flex flex-col gap-2.5">
        <Label.Root className="text-2xl" htmlFor={templateEditorId}>
          Template:
        </Label.Root>
        <CodeEditor id={templateEditorId} value={template} onValueChange={(v) => setTemplate(v)} />
      </div>

      <div className="flex flex-col gap-2.5">
        <Label.Root className="text-2xl">
          <div className="flex gap-2 justify-between place-items-center">
            <div>Result:</div>
            <div>{updating && <LoadingIcon size="1.25rem" />}</div>
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
          disabled={updating || isError}
          onClick={() => updateField(template, renderedField?.result)}
          className="flex-1"
        >
          Update field
        </Button>
        <Button
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
  example?: string;
  error?: string;
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
  if (!resp.Anki.anki) {
    return [[], []];
  }
  const fields: Array<MappingField> = resp.Anki.anki.noteFields.map((e) => {
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
