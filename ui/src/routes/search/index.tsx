import { clsx } from 'clsx';
import { useEffect, useRef } from 'react';
import {
  Form,
  LoaderFunctionArgs,
  useLoaderData,
  useNavigate,
  useNavigation,
} from 'react-router-dom';

import { gql } from '../../api/__generated__/gql';
import { LemmaNoteInfo } from '../../api/__generated__/graphql';
import apolloClient from '../../apollo-client';
import LensIcon from '../../components/Icons/LensIcon';
import { LoaderData } from '../../loader-type';
import LemmaList from './LemmaList';

const GET_LEMMAS = gql(`
  query GetLemmas($query: String!) {
    Lemmas(query: $query) {
      lemmas {
        noteID
        lemma {
          slug {
            word
            hiragana
            furigana {
              kanji
              hiragana
            }
            pitchShapes {
              hiragana
              directions
            } 
          }
          tags
          forms {
            word
            hiragana
            furigana {
              kanji
              hiragana
            }
            pitchShapes {
              hiragana
              directions
            } 
          }
          definitions
          partsOfSpeech
          senseTags
          audio {
            mediaType
            source
          }
        } 
      }
    }
  }
`);

export async function loader({ params }: LoaderFunctionArgs): Promise<{
  lemmas: LemmaNoteInfo[] | undefined;
  query: string;
}> {
  const query = params.query ?? '';
  if (!query) {
    return {
      lemmas: undefined,
      query,
    };
  }
  const { data } = await apolloClient.query({
    query: GET_LEMMAS,
    variables: {
      query: query,
    },
  });

  return { lemmas: data.Lemmas?.lemmas, query };
}

enum SearchState {
  Typing,
  Routing,
}

export default function Search() {
  const { lemmas, query } = useLoaderData() as LoaderData<typeof loader>;
  const navigate = useNavigate();
  const navigation = useNavigation();
  const inputQuery = useRef<HTMLInputElement>(null);
  // searchState needed to track what user doing: because we want to change inputQuery only in case
  // when user changed query manually or typed submit and did type anything in form
  const searchState = useRef<SearchState>(SearchState.Typing);
  const handleChange = () => {
    searchState.current = SearchState.Typing;
  };
  const handleSubmit: React.FormEventHandler = (event) => {
    event.preventDefault();
    searchState.current = SearchState.Routing;
    if (!inputQuery.current) {
      return;
    }
    // TODO: this is ugly, but workarounds seems to be complicated
    navigate(`/search/${inputQuery.current.value}`);
  };
  useEffect(() => {
    if (inputQuery.current && searchState.current == SearchState.Routing) {
      inputQuery.current.value = query;
    }
  }, [query]);
  useEffect(() => {
    searchState.current = SearchState.Routing;
  }, [navigation.state]);

  return (
    <div className="flex flex-1 flex-col">
      <Form className="group relative" onSubmit={handleSubmit}>
        <LensIcon className="text-slate-400 group-focus-within:text-blue-500 pointer-events-none absolute left-3 top-1/2 -mt-2.5 h-5 w-5" />
        <input
          id="query"
          className="text-slate-900 placeholder-slate-400 w-full appearance-none rounded-md py-2 pl-10 text-lg leading-6 shadow-sm ring-1 ring-blue focus:outline-none focus:ring-2"
          type="text"
          aria-label="Enter japanese word"
          placeholder="Enter japanese word..."
          defaultValue={query}
          ref={inputQuery}
          onChange={handleChange}
        />
      </Form>
      <div
        className={clsx(
          'flex flex-1',
          navigation.state == 'loading' && 'opacity-30 transition-opacity delay-300 duration-300',
        )}
      >
        {lemmas && lemmas.length != 0 ? (
          <div className="flex flex-1 flex-col">
            <LemmaList key={query} lemmaNotes={lemmas} />
          </div>
        ) : (
          <div className="flex flex-1 flex-col justify-center">
            <div className="flex justify-center">
              <h1 className="text-5xl text-blue">
                {lemmas === undefined ? 'Make your search request' : 'Nothing was found'}
              </h1>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
