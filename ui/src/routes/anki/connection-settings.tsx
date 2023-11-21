import { useMutation, useSuspenseQuery } from '@apollo/client';
import * as Label from '@radix-ui/react-label';
import { clsx } from 'clsx';
import { ErrorMessage, Form, Formik, FormikErrors, FormikHelpers, FormikProps } from 'formik';
import { useCallback, useContext, useEffect, useId, useRef } from 'react';
import { Link } from 'react-router-dom';
import * as yup from 'yup';

import { gql } from '../../api/__generated__/gql';
import Button, { ButtonVariant } from '../../components/Button';
import StatusIcon, { StatusIconKind } from '../../components/StatusIcon';
import SuspenseLoading from '../../components/SuspenseLoading';
import TextField, { TextFieldProps } from '../../components/TextField';
import { HealthStatusContext, HealthStatusReloadContext } from '../../contexts/health-status';
import { useToastify } from '../../hooks/toastify';
import { apolloErrorToast } from '../../lib/styled-toast';
import { validateHostPort } from '../../lib/validate';
import {
  AnkiState,
  HealthStatus,
  HealthStatusLoading,
  HealthStatusOk,
  throwErrorHealthStatus,
} from '../../model/health-status';

export default function AnkiConnectionSettings() {
  return (
    <div className="flex flex-col gap-8">
      <div>
        <p className="text-blue text-lg">
          <span className="font-bold">Note:</span>
          <br />
          To use this application with{' '}
          <a href="https://apps.ankiweb.net/" target="_blank" rel="noreferrer">
            <span className="text-dark-blue">Anki</span>
          </a>{' '}
          you should install{' '}
          <a href="https://ankiweb.net/shared/info/2055492159" target="_blank" rel="noreferrer">
            <span className="text-dark-blue">AnkiConnect</span>
          </a>{' '}
          plugin. After you need to start Anki and enter your account.
        </p>
      </div>
      <SuspenseLoading>
        <FormConnectionSettings />
      </SuspenseLoading>
    </div>
  );
}

const GET_CONNECTION_CONFIG = gql(`
  query GetConnectionConfig {
    AnkiConfig {
      addr
      apiKey
    }
  }
`);

const UPDATE_CONNECTION_CONFIG = gql(`
  mutation UpdateConnectionConfig($addr: String!, $apiKey: String!) {
    setAnkiConfigConnection(input: { addr: $addr, apiKey: $apiKey }) {
      error {
          ... on ValidationError {
            paths
            message
          }
        }
      }
  }
`);

type ConnectionSettingsData = {
  addr: string;
  apiKey: string;
};

const ConnectionSettingsValidationSchema = yup.object({
  addr: yup
    .string()
    .required()
    .test((value, context) => {
      const error = validateHostPort(value);
      if (error) {
        return context.createError({ message: error });
      }
      return true;
    }),
});

function FormConnectionSettings() {
  const healthStatus = useContext(HealthStatusContext);
  throwErrorHealthStatus(healthStatus);
  const healthStatusReload = useContext(HealthStatusReloadContext);
  const [updateConnectionConfig] = useMutation(UPDATE_CONNECTION_CONFIG, {
    update: (cache) => {
      cache.evict({ id: 'AnkiConfig:{}' });
    },
    // we passe empty onError because in this case apollo will not throw in case if network error for some reason
    onError: () => {},
  });
  const { data } = useSuspenseQuery(GET_CONNECTION_CONFIG);
  const { addr, apiKey } = data.AnkiConfig;
  const toast = useToastify({
    type: 'success',
  });
  const handleSubmit = useCallback(
    async (values: ConnectionSettingsData, helper: FormikHelpers<ConnectionSettingsData>) => {
      const { data, errors } = await updateConnectionConfig({ variables: values });
      if (errors) {
        apolloErrorToast(errors, 'Anki Connect settings change failed.', { toast: toast });
        healthStatusReload();
        return;
      }
      const response = data?.setAnkiConfigConnection.error;
      if (response?.__typename == 'ValidationError') {
        helper.setErrors({ addr: response.message });
        toast('Anki Connect settings change failed.', { type: 'error' });
      } else {
        helper.resetForm({ values: values });
        toast('Anki Connect settings saved.');
        healthStatusReload();
      }
    },
    [toast, updateConnectionConfig, healthStatusReload],
  );
  const formRef = useRef<FormikProps<ConnectionSettingsData>>(null);
  useEffect(() => {
    if (!formRef.current || healthStatus.kind == 'Loading') {
      return;
    }
    const statusErrors = userErrorsFromHealthStatus(healthStatus);
    if (ConnectionSettingsValidationSchema.isValidSync(formRef.current.values)) {
      formRef.current.setErrors(statusErrors);
    }
  }, [healthStatus]);

  return (
    <>
      <Formik<ConnectionSettingsData>
        initialValues={{ addr: addr, apiKey: apiKey }}
        validationSchema={ConnectionSettingsValidationSchema}
        onSubmit={handleSubmit}
        innerRef={formRef}
        initialTouched={{ addr: true, apiKey: true }}
      >
        {(props) => {
          return (
            <>
              <Form className="flex flex-col gap-8">
                <FormTextField
                  name="addr"
                  type="input"
                  label="Anki Address"
                  inputClassName="max-w-md"
                />
                <FormTextField
                  name="apiKey"
                  type="password"
                  autoComplete="off"
                  label="API Key"
                  inputClassName="max-w-md"
                />
                <div className="flex flex-row gap-8">
                  <Button type="submit" className="basis-52" disabled={props.isSubmitting}>
                    Save Settings
                  </Button>
                  <Button
                    type="button"
                    onClick={props.handleReset}
                    className={clsx('basis-52')}
                    variant={ButtonVariant.Dangerous}
                    disabled={props.isSubmitting}
                  >
                    Reset
                  </Button>
                </div>
              </Form>
              <StatusBox status={healthStatus} />
            </>
          );
        }}
      </Formik>
    </>
  );
}

type FormTextFieldProps = TextFieldProps & { label: string };

function FormTextField({ name, label, ...rest }: FormTextFieldProps) {
  const inputId = useId();
  return (
    <div className={clsx('flex flex-col gap-2.5 text-2xl')}>
      {label && <Label.Root htmlFor={inputId}>{label}</Label.Root>}
      <div className="flex  flex-col  gap-2.5">
        <TextField id={inputId} name={name} {...rest} />
        <p className="text-error-red text-lg">
          <ErrorMessage name={name} />
        </p>
      </div>
    </div>
  );
}

function StatusBox({ status }: { status: HealthStatusOk | HealthStatusLoading }) {
  const { iconKind, title, titleClassName, body } = getStatusBoxContent(status);
  return (
    <div className="flex flex-row justify-start items-start gap-2 basis-16">
      <StatusIcon size="2.5rem" kind={iconKind} />
      <div className="text-2xl">
        <h1 className={clsx('text-bold leading-10', titleClassName)}>{title}</h1>
        {body}
      </div>
    </div>
  );
}

function getStatusBoxContent(status: HealthStatusOk | HealthStatusLoading): StatusBoxContent {
  switch (status.kind) {
    case 'Loading':
      return {
        iconKind: 'Loading',
        title: 'Loading...',
        titleClassName: 'text-dark-blue',
        body: null,
      };
    case 'Ok':
      return kindToStatusBoxContent[status.anki.kind];
    default: {
      const _exhaustiveCheck: never = status;
      return _exhaustiveCheck;
    }
  }
}

type StatusBoxContent = {
  iconKind: StatusIconKind;
  title: string;
  titleClassName: string;
  body: React.ReactNode;
};

const kindToStatusBoxContent: Record<AnkiState['kind'], StatusBoxContent> = {
  Ok: {
    iconKind: 'OK',
    title: 'Connected',
    titleClassName: 'text-green',
    body: <p>All is OK!</p>,
  },
  UserError: {
    iconKind: 'OK',
    title: 'Connected',
    titleClassName: 'text-green',
    body: (
      <p>
        Go to{' '}
        <Link to="../user-settings" className="text-blue underline">
          Anki User Setting
        </Link>{' '}
        to complete connection
      </p>
    ),
  },
  ConnectionError: {
    iconKind: 'Error',
    title: 'Error',
    titleClassName: 'text-light-red',
    body: <p>Unable to connect to AnkiConnect with specified address</p>,
  },
  ForbiddenOrigin: {
    iconKind: 'Error',
    title: 'Error',
    titleClassName: 'text-light-red',
    body: <p>Origin is unknown</p>,
  },
  Unauthorized: {
    iconKind: 'Error',
    title: 'Error',
    titleClassName: 'text-light-red',
    body: <p>Invalid API Key</p>,
  },
};

function userErrorsFromHealthStatus(status: HealthStatus): FormikErrors<ConnectionSettingsData> {
  switch (status.kind) {
    case 'Ok':
      switch (status.anki.kind) {
        case 'ConnectionError':
          return {
            addr: `unable to connect: ${status.anki.message}`,
          };
        case 'ForbiddenOrigin':
          return {
            addr: 'add client origin to anki-connect whitelist',
          };
        case 'Unauthorized':
          return {
            apiKey: 'API Key is incorrect',
          };
      }
  }
  return {};
}
