import {
  Form,
  LoaderFunctionArgs,
  useNavigate,
  useLoaderData,
  useNavigation,
} from 'react-router-dom';
import LemmaList from '../components/LemmaList';
import { LoaderData } from '../loader-type';
import apolloClient from '../apollo-client';
import { Lemma } from '../api/__generated__/graphql';
import { gql } from '../api/__generated__/gql';
import { clsx } from 'clsx';
import { useContext, useEffect } from 'react';
import { RouterContext } from '../context/RouterContext';

const GET_LEMMAS = gql(`
  query GetLemmas($query: String!) {
    Lemmas(query: $query) {
      lemmas {
        slug {
          word
          hiragana
          furigana {
            kanji
            hiragana
          }
          pitch {
            hiragana
            pitch
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
          pitch {
            hiragana
            pitch
          } 
        }
        senses {
          definition
          partOfSpeech
          tags
        }
        audio {
          type
          source
        }
      } 
    }
  }
`);

export async function loader({ params }: LoaderFunctionArgs): Promise<{
  lemmas: Array<Lemma> | undefined;
  query: string;
}> {
  const query = params.query ?? '';
  if (!query) {
    return {
      lemmas: undefined,
      query,
    };
  }
  const lemmas = await apolloClient.query({
    query: GET_LEMMAS,
    variables: {
      query: query,
    },
  });

  return { lemmas: lemmas.data.Lemmas?.lemmas, query };
}

export default function Search() {
  const { lemmas, query } = useLoaderData() as LoaderData<typeof loader>;
  const navigate = useNavigate();
  const navigation = useNavigation();
  const handleSubmit: React.FormEventHandler = (event) => {
    event.preventDefault();
    const queryInput = document.getElementById('query') as HTMLInputElement;
    // TODO: this is ugly, but workarounds seems to be complicated
    navigate(`/search/${queryInput.value}`);
  };
  useEffect(() => {
    (document.getElementById('query') as HTMLInputElement).value = query;
  }, [query]);

  return (
    <>
      <Form className="group relative" onSubmit={handleSubmit}>
        <svg
          width="20"
          height="20"
          fill="currentColor"
          className="absolute left-3 top-1/2 -mt-2.5 text-slate-400 pointer-events-none group-focus-within:text-blue-500"
          aria-hidden="true"
        >
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
          />
        </svg>
        <input
          id="query"
          className="focus:ring-2 focus:ring-green focus:outline-none appearance-none w-full text-sm leading-6 text-slate-900 placeholder-slate-400 rounded-md py-2 pl-10 ring-1 ring-blue shadow-sm"
          type="text"
          aria-label="Enter japanese word"
          placeholder="Enter japanese word..."
          defaultValue={query}
        />
      </Form>
      <div
        className={clsx(
          'flex-1 flex flex-col justify-center',
          navigation.state == 'loading' && 'opacity-30 transition-opacity duration-300 delay-300',
        )}
      >
        {lemmas && lemmas.length != 0 ? (
          <LemmaList lemmas={lemmas} />
        ) : (
          <div className="flex justify-center">
            <h1 className="text-5xl text-blue">
              {lemmas === undefined ? 'Make your search request' : 'Nothing was found'}
            </h1>
          </div>
        )}
      </div>
    </>
  );
}
