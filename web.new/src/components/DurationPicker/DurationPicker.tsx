import { Input } from '../Input';
import styled from 'styled-components';
import { useRef } from 'react';
import { useTranslation } from 'react-i18next';

type Props = React.HTMLAttributes<HTMLDivElement> & {
  value?: number;
  onDurationInput?: (v: number) => void;
};

const PickerCotnainer = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 1em;
  ${Input} {
    margin-top: 0.5em;
    display: block;
    background-color: ${(p) => p.theme.background3};
    width: calc(6ch + 20px);
  }
`;

export const DurationPicker: React.FC<Props> = ({
  value = 0,
  onDurationInput = () => {},
  ...props
}) => {
  const { t } = useTranslation('components', { keyPrefix: 'durationpicker' });

  const daysRef = useRef<HTMLInputElement>(null);
  const hoursRef = useRef<HTMLInputElement>(null);
  const minutesRef = useRef<HTMLInputElement>(null);
  const secondsRef = useRef<HTMLInputElement>(null);

  const _onInput = () => {
    console.log(
      combineSeconds(
        parseInt(daysRef.current?.value || '0'),
        parseInt(hoursRef.current?.value || '0'),
        parseInt(minutesRef.current?.value || '0'),
        parseInt(secondsRef.current?.value || '0'),
      ),
    );
    onDurationInput(
      combineSeconds(
        parseInt(daysRef.current?.value || '0'),
        parseInt(hoursRef.current?.value || '0'),
        parseInt(minutesRef.current?.value || '0'),
        parseInt(secondsRef.current?.value || '0'),
      ),
    );
  };

  const [days, hours, minutes, seconds] = spreadSeconds(value);

  return (
    <PickerCotnainer {...props}>
      <section>
        <label htmlFor=":duration-days">{t('days')}</label>
        <Input
          ref={daysRef}
          id=":duration-days"
          type="number"
          min={0}
          value={days}
          onInput={_onInput}
        />
      </section>
      <section>
        <label htmlFor=":duration-hours">{t('hours')}</label>
        <Input
          ref={hoursRef}
          id=":duration-hours"
          type="number"
          min={0}
          max={23}
          value={hours}
          onInput={_onInput}
        />
      </section>
      <section>
        <label htmlFor=":duration-minutes">{t('minutes')}</label>
        <Input
          ref={minutesRef}
          id=":duration-minutes"
          type="number"
          min={0}
          max={59}
          value={minutes}
          onInput={_onInput}
        />
      </section>
      <section>
        <label htmlFor=":duration-seconds">{t('seconds')}</label>
        <Input
          ref={secondsRef}
          id=":duration-seconds"
          type="number"
          min={0}
          max={59}
          value={seconds}
          onInput={_onInput}
        />
      </section>
    </PickerCotnainer>
  );
};

const spreadSeconds = (value: number): [number, number, number, number] => {
  const days = Math.floor(value / (24 * 3600));
  const hours = Math.floor((value % (24 * 3600)) / 3600);
  const minutes = Math.floor(((value % (24 * 3600)) % 3600) / 60);
  const seconds = Math.floor(((value % (24 * 3600)) % 3600) % 60);
  return [days, hours, minutes, seconds];
};

const combineSeconds = (days: number, hours: number, minutes: number, seconds: number): number =>
  days * 24 * 3600 + hours * 3600 + minutes * 60 + seconds;
