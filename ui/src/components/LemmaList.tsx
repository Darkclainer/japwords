import { Lemma } from '../api/__generated__/graphql';
import LemmaCard from './LemmaCard';

export default function LemmaList(props: { lemmas: Array<Lemma> }) {
  const { lemmas } = props;
  return (
    <>
      {lemmas.map((lemma, index) => (
        <LemmaCard key={lemma.slug.word + '-' + index} lemma={lemma} />
      ))}
    </>
  );
}
