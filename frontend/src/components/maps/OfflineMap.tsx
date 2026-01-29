import { useEffect, useRef } from 'react';
import maplibregl from 'maplibre-gl';
import { Protocol } from 'pmtiles';

export default function OfflineMap() {
  const mapRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (!mapRef.current) return;

    // Регистрируем pmtiles:// протокол
    const protocol = new Protocol();
    maplibregl.addProtocol('pmtiles', protocol.tile);

    const map = new maplibregl.Map({
      container: mapRef.current,
      style: {
        version: 8,
        sources: {
          openmaptiles: {
            type: 'vector',
            url: 'pmtiles:///maps/turkmenistan.pmtiles',
          },
        },
        layers: [
          {
            id: 'water',
            type: 'fill',
            source: 'openmaptiles',
            'source-layer': 'water',
            paint: { 'fill-color': '#a0c8f0' },
          },
          {
            id: 'roads',
            type: 'line',
            source: 'openmaptiles',
            'source-layer': 'transportation',
            paint: {
              'line-color': '#ffffff',
              'line-width': 1,
            },
          },
          {
            id: 'buildings',
            type: 'fill',
            source: 'openmaptiles',
            'source-layer': 'building',
            paint: {
              'fill-color': '#d0d0d0',
              'fill-opacity': 0.6,
            },
          },
        ],
      },
      center: [59.6, 38.9], // Туркменистан
      zoom: 5,
    });

    return () => map.remove();
  }, []);

  return <div ref={mapRef} className="w-full h-[600px] rounded-lg border" />;
}
