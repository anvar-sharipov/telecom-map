import Button from '../../components/UI/Button/Button';

type MapMode = 'idle' | 'create-telecom' | 'create-magistral';

type MainButtonsProps = {
  zoomToDashoguz: () => void;
  mode: MapMode;
  setMode: (mode: MapMode) => void;
};
const MainButtons = ({ zoomToDashoguz, mode, setMode }: MainButtonsProps) => {
  return (
    <div>
      <div
        style={{
          position: 'absolute',
          zIndex: 10,
          top: 130,
          left: 20,
          display: 'flex',
          flexDirection: 'column',
          gap: 10,
        }}
      >
        <Button onClick={zoomToDashoguz}>Dashoguz</Button>

        {/* 🔐 только для админа */}
        <Button
          variant={mode === 'create-telecom' ? 'secondary' : 'danger'}
          onClick={() => setMode('create-telecom')}
        >
          ➕ Создать Telecom
        </Button>
      </div>
    </div>
  );
};

export default MainButtons;
