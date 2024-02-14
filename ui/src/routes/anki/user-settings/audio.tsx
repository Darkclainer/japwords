import * as Label from '@radix-ui/react-label';
import * as Select from '@radix-ui/react-select';
import {
  SelectContent,
  SelectItem,
  SelectScrollDownButton,
  SelectScrollUpButton,
  SelectSeparator,
  SelectTrigger,
} from '../../../components/Select';
import { useCallback, useMemo, useState } from 'react';
import { gql } from '../../../api/__generated__';
import { useMutation } from '@apollo/client';
import { GET_HEALTH_STATUS } from '../../../api/health-status';
import { GET_ANKI_CONFIG } from './api';
import { useToastify } from '../../../hooks/toastify';
import { apolloErrorToast } from '../../../lib/styled-toast';
import TextInput from '../../../components/TextInput';

type AudioEditProps = {
  audioField: string;
  audioPreferredType: string;
  ankiNoteFields: string[];
};

export default function AudioEdit(props: AudioEditProps) {
  return (
    <>
      <div className="flex flex-col gap-2.5">
        <Label.Root className="text-2xl">Audio mapping field:</Label.Root>
        <AudioFieldSelect {...props} />
      </div>
      <div className="flex flex-col gap-2.5">
        <Label.Root className="text-2xl">Audio preferred type:</Label.Root>
        <AudioPreferredTypeInput {...props} />
      </div>
    </>
  );
}

const SET_AUDIO_FIELD = gql(`
  mutation SetAnkiConfigAudioField($field: String!) {
    setAnkiConfigAudioField(input: { audioField: $field }) {
      error {
          message
      }
    }
  }
`);

// quote can not be used in field names
const audioFieldDisabledValue = '"__disabled"';

function AudioFieldSelect({
  audioField,
  ankiNoteFields,
}: {
  audioField: string;
  ankiNoteFields: string[];
}) {
  const [setAudioField, { loading }] = useMutation(SET_AUDIO_FIELD, {
    refetchQueries: [GET_HEALTH_STATUS, GET_ANKI_CONFIG],
    awaitRefetchQueries: true,
  });
  const toast = useToastify({ type: 'success' });
  const onValueChange = useCallback(
    async (value: string) => {
      if (value === audioFieldDisabledValue) {
        value = '';
      }
      const { data, errors } = await setAudioField({
        variables: {
          field: value,
        },
      });
      if (errors) {
        apolloErrorToast(errors, 'Audio mapping field update failed.', { toast: toast });
        return;
      } else if (data?.setAnkiConfigAudioField.error) {
        toast('Audio mapping field update failed. ' + data.setAnkiConfigAudioField.error.message, {
          type: 'error',
        });
        return;
      }
      toast('Audio mapping field updated.');
    },
    [setAudioField, toast],
  );
  const fieldNotExists = useMemo(
    () => audioField != '' && ankiNoteFields.indexOf(audioField) === -1,
    [audioField, ankiNoteFields],
  );
  return (
    <>
      <Select.Root onValueChange={onValueChange} value={audioField} disabled={loading}>
        <SelectTrigger className="max-w-md" hasError={fieldNotExists}>
          <Select.Value placeholder={'Audio mapping disabled'} asChild>
            <span>{audioField}</span>
          </Select.Value>
        </SelectTrigger>
        <Select.Portal>
          <SelectContent>
            <Select.Group className="group-data-side-top/content:order-3">
              <SelectItem value={audioFieldDisabledValue}>
                <Select.ItemText>Disable audio mapping</Select.ItemText>
              </SelectItem>
            </Select.Group>
            <SelectSeparator className="group-data-side-top/content:order-2" />
            <SelectScrollUpButton />
            <Select.Viewport>
              <Select.Group>
                {ankiNoteFields.map((field, i) => (
                  <SelectItem key={i} value={field}>
                    <Select.ItemText>{field}</Select.ItemText>
                  </SelectItem>
                ))}
              </Select.Group>
            </Select.Viewport>
            <SelectScrollDownButton />
          </SelectContent>
        </Select.Portal>
      </Select.Root>
      {fieldNotExists && (
        <p className="text-lg text-error-red">Selected field does not exists in the note type</p>
      )}
    </>
  );
}

const SET_AUDIO_PREFERRED_TYPE = gql(`
  mutation SetAnkiConfigAudioPreferredType($preferredType: String!) {
    setAnkiConfigAudioPreferredType(input: { audioPreferredType: $preferredType }) {
      nothing
    }
  }
`);

function AudioPreferredTypeInput({ audioPreferredType }: { audioPreferredType: string }) {
  const [setAudioPreferredType, { loading }] = useMutation(SET_AUDIO_PREFERRED_TYPE, {
    refetchQueries: [GET_HEALTH_STATUS, GET_ANKI_CONFIG],
    awaitRefetchQueries: true,
  });
  const [preferredType, setPreferredType] = useState(audioPreferredType);
  return (
    <TextInput
      className="max-w-md"
      value={preferredType}
      onChange={(e) => setPreferredType(e.target.value)}
    />
  );
}
