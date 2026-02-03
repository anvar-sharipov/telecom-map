import React, { useEffect, useRef, useState } from 'react';
import maplibregl from 'maplibre-gl';
import { Protocol } from 'pmtiles';

import 'maplibre-gl/dist/maplibre-gl.css';
import MainButtons from './map_pages/MainButtons';

import { icons } from './map_pages/icons';
import type { FeatureCollection, Feature, Point, LineString } from 'geojson';

const Home: React.FC = () => {
  const mapRef = useRef<HTMLDivElement | null>(null);
  const mapInstanceRef = useRef<maplibregl.Map | null>(null);

  const pointSourceRef = useRef<FeatureCollection<Point>>({
    type: 'FeatureCollection',
    features: [],
  });

  const lineSourceRef = useRef<FeatureCollection<LineString>>({
    type: 'FeatureCollection',
    features: [],
  });

  type MapMode = 'idle' | 'create-telecom' | 'create-magistral';
  const [mode, setMode] = useState<MapMode>('idle');
  const modeRef = useRef<MapMode>('idle');

  useEffect(() => {
    modeRef.current = mode;
  }, [mode]);

  useEffect(() => {
    if (!mapRef.current) return;

    const protocol = new Protocol();
    maplibregl.addProtocol('pmtiles', protocol.tile);

    const map = new maplibregl.Map({
      container: mapRef.current,
      style: '/maps/style.json',
      center: [59.6, 39.0],
      zoom: 5,
    });

    mapInstanceRef.current = map;

    map.on('load', async () => {
      for (const name of icons) {
        try {
          const image = await map.loadImage(`/maps/icons/${name}.png`);

          if (!map.hasImage(name)) {
            map.addImage(name, image.data, { sdf: true });
          }
        } catch (e) {
          console.error('❌ icon load error:', name, e);
        }
      }

      map.addSource('click-point', {
        type: 'geojson',
        data: pointSourceRef.current,
      });

      map.addSource('telecom-line', {
        type: 'geojson',
        data: lineSourceRef.current,
      });

      map.addLayer({
        id: 'click-point-layer',
        type: 'circle',
        source: 'click-point',
        paint: {
          'circle-radius': 6,
          'circle-color': 'green',
        },
      });

      map.addLayer({
        id: 'telecom-line-layer',
        type: 'line',
        source: 'telecom-line',
        paint: {
          'line-color': '#0066ff',
          'line-width': 3,
        },
      });

      map.addSource('telecom', {
        type: 'geojson',
        data: '/data/telecom_objects.geojson',
      });

      // ✅ 3. Добавляем слой ПОСЛЕ иконок
      map.addLayer({
        id: 'telecom-icons',
        type: 'symbol',
        source: 'telecom',
        minzoom: 14,
        layout: {
          'icon-image': [
            'match',
            ['get', 'type'],

            'cabinet',
            'cabinet',
            'splice',
            'splice',
            'home',
            'home',

            'marker', // fallback
          ],
          'icon-size': 1,
          'icon-allow-overlap': true,
        },
      });
    });

    map.on('click', (e) => {
      if (modeRef.current !== 'create-telecom') {
        return;
      }

      const { lng, lat } = e.lngLat;

      pointSourceRef.current.features.push({
        type: 'Feature',
        geometry: {
          type: 'Point',
          coordinates: [lng, lat],
        },
        properties: { type: 'telecom-node' },
      });

      (map.getSource('click-point') as maplibregl.GeoJSONSource).setData(pointSourceRef.current);

      lineSourceRef.current.features = [
        {
          type: 'Feature',
          geometry: {
            type: 'LineString',
            coordinates: pointSourceRef.current.features.map((f) => f.geometry.coordinates),
          },
          properties: {},
        },
      ];

      (map.getSource('telecom-line') as maplibregl.GeoJSONSource).setData(lineSourceRef.current);
      // if (modeRef.current === 'create-telecom') {
      //   console.log('📡 CREATE TELECOM:', lng, lat);

      //   // ⏳ дальше тут будет:
      //   // add point to geojson
      //   // or POST /nodes

      //   setMode('idle');
      //   return;
      // }

      // pointSource.features = [
      //   {
      //     type: 'Feature',
      //     geometry: {
      //       type: 'Point',
      //       coordinates: [lng, lat],
      //     },
      //     properties: {},
      //   },
      // ];
      // (map.getSource('click-point') as maplibregl.GeoJSONSource).setData(pointSource);

      // pointSource.features.push({
      //   type: 'Feature',
      //   geometry: {
      //     type: 'Point',
      //     coordinates: [lng, lat],
      //   },
      //   properties: {},
      // });
      // (map.getSource('click-point') as maplibregl.GeoJSONSource).setData(pointSource);
    });

    (window as any).map = map;

    return () => {
      map.remove();
      maplibregl.removeProtocol('pmtiles');
    };
  }, []);

  const undoLastPoint = () => {
    const map = mapInstanceRef.current;
    if (!map) return;

    const features = pointSourceRef.current.features;
    if (features.length === 0) return;

    features.pop();

    (map.getSource('click-point') as maplibregl.GeoJSONSource).setData(pointSourceRef.current);

    if (features.length >= 2) {
      lineSourceRef.current.features = [
        {
          type: 'Feature',
          geometry: {
            type: 'LineString',
            coordinates: features.map((f) => f.geometry.coordinates),
          },
          properties: {},
        },
      ];
    } else {
      lineSourceRef.current.features = [];
    }

    (map.getSource('telecom-line') as maplibregl.GeoJSONSource).setData(lineSourceRef.current);
  };

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if ((e.ctrlKey && e.key === 'z') || e.key === 'Backspace') {
        undoLastPoint();
      }
    };

    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  }, []);

  const zoomToDashoguz = () => {
    if (!(window as any).map) return;

    (window as any).map.fitBounds(
      [
        [58.5, 41.1], // юго-запад
        [60.4, 42.5], // северо-восток
      ],
      {
        padding: {
          top: 120, // 👈 под кнопку
          bottom: 40, // 👈 показываем низ
          left: 30,
          right: 30,
        },
        duration: 1200,
      },
    );
  };
  return (
    <>
      <div ref={mapRef} style={{ width: '100%', height: '100vh' }} />

      <div>
        <MainButtons zoomToDashoguz={zoomToDashoguz} mode={mode} setMode={setMode} />
      </div>
    </>
  );
};

export default Home;
