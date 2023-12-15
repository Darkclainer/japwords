import { useLazyQuery } from '@apollo/client';
import { useCallback, useMemo, useRef, useState } from 'react';

import { gql } from '../../api/__generated__';
import { LemmaNoteInfo } from '../../api/__generated__/graphql';
import { useToastify } from '../../hooks/toastify';
import { groupLemmaNotes } from '../../lib/lemma-bag';
import { apolloErrorToast } from '../../lib/styled-toast';
import LemmaCard from './LemmaCard';
import AddNoteDialog, { AddLemmaRequest } from './AddNoteDialog';

const PREPARE_LEMMA = gql(`
query PrepareLemma($lemma: LemmaInput) {
  PrepareLemma(lemma: $lemma) {
    request {
      fields {
        name
        value
      }
      tags
      audioAssets {
        field
        filename
        url
        data
      }
    }
    error {
      ... on AnkiIncompleteConfiguration {
        message
      }
      ... on Error {
        message
      }
    }
    ankiError {
      __typename
      ... on Error {
        message
      }
    }
  }
}
`);

export default function LemmaList({
  lemmaNotes: initialLemmaNotes,
}: {
  lemmaNotes: LemmaNoteInfo[];
}) {
  // more propper way will be to modify cache, but because lemmas now returned from react-router-dom it's not possible
  const [lemmaNotes, setLemmaNotes] = useState(initialLemmaNotes);
  const lemmaBags = useMemo(() => groupLemmaNotes(lemmaNotes), [lemmaNotes]);
  const [prepareLemma] = useLazyQuery(PREPARE_LEMMA, {
    fetchPolicy: 'network-only',
  });
  const [openAddNote, setOpenAddNote] = useState(false);
  const [addLemmaRequest, setAddLemmaRequest] = useState<AddLemmaRequest>();
  const abortReset = useAbortRef();
  const toast = useToastify({
    autoClose: 2000,
    type: 'info',
  });
  const previewLemma = useCallback(
    async (lemmaNote: LemmaNoteInfo) => {
      const currentAbortController = abortReset();
      setAddLemmaRequest(undefined);
      setOpenAddNote(true);
      const now = new Date().valueOf();
      const { data, error } = await prepareLemma({
        variables: {
          lemma: lemmaNote.lemma,
        },
        context: {
          fetchOptions: {
            signal: currentAbortController.signal,
          },
        },
      });
      // if request is to fast, display loading for minimum 400ms
      // Better, to show window if request is longer then 200ms for example?
      const waitTime = Math.max(0, -(new Date().valueOf() - now) + 400);
      if (waitTime > 0) {
        await new Promise((resolve) => setTimeout(resolve, waitTime));
      }
      if (error) {
        setOpenAddNote(false);
        apolloErrorToast(error, `${AddNoteFailedActionTitle}: `, {
          toast,
        });
        return;
      }
      const ankiError = data?.PrepareLemma.ankiError;
      if (ankiError) {
        setOpenAddNote(false);
        toast(`${AddNoteFailedActionTitle}: problems with anki-connect`, {
          type: 'error',
        });
        return;
      }
      const userError = data?.PrepareLemma.error;
      if (userError) {
        setOpenAddNote(false);
        let message = `${AddNoteFailedActionTitle}: `;
        if (userError.__typename == 'AnkiIncompleteConfiguration') {
          message += 'anki configuration unfinished';
        } else {
          message += 'uknown reason';
        }
        toast(message, {
          type: 'error',
        });
        return;
      }
      if (!currentAbortController.signal.aborted) {
        if (!data?.PrepareLemma.request) {
          throw 'unreachable';
        }
        setAddLemmaRequest({
          note: data.PrepareLemma.request,
          lemma: lemmaNote,
        });
      }
    },
    [setOpenAddNote, toast, prepareLemma],
  );

  return (
    <>
      {lemmaBags.map((lemmaBag, index) => (
        <LemmaCard
          key={lemmaBag.slug.word + '-' + index}
          lemmaBag={lemmaBag}
          toast={toast}
          previewLemma={previewLemma}
        />
      ))}
      <AddNoteDialog
        open={openAddNote}
        setOpen={(value: boolean) => {
          if (!value) {
            abortReset();
          }
          return setOpenAddNote(value);
        }}
        addLemmaRequest={addLemmaRequest}
        toast={toast}
        setLemmaNotes={setLemmaNotes}
      />
    </>
  );
}

function useAbortRef() {
  const abortRef = useRef<AbortController>();
  const abortReset = useCallback(() => {
    const previousController = abortRef.current;
    if (previousController) {
      previousController.abort();
    }
    const newController = new AbortController();
    abortRef.current = newController;
    return newController;
  }, []);
  return abortReset;
}
