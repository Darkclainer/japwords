import { Lemma } from '../api/__generated__/graphql';
import { useToastify } from '../hooks/toastify';
import LemmaCard from './LemmaCard';

export default function LemmaList(props: { lemmas: Array<Lemma> }) {
  const { lemmas } = props;
  const toast = useToastify({
    position: 'bottom-right',
    autoClose: 2000,
    type: 'info',
  });
  return (
    <>
      {lemmas.map((lemma, index) => (
        <LemmaCard key={lemma.slug.word + '-' + index} lemma={lemma} toast={toast} />
      ))}
    </>
  );
}
