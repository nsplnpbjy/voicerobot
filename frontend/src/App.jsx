import React, { useState, useEffect, useRef } from 'react';
import { V2T, T2V, VMT } from "../wailsjs/go/main/App";

const App = () => {
  const [isRecording, setIsRecording] = useState(false);
  const [audioChunks, setAudioChunks] = useState([]);
  const [audioUrl, setAudioUrl] = useState(null);
  const [audioDownloadUrl, setAudioDownloadUrl] = useState(null);
  const [mediaRecorder, setMediaRecorder] = useState(null);
  const [transcription, setTranscription] = useState("");
  const [vmtResult, setVmtResult] = useState("");
  const audioRef = useRef(null);
  const [audioSrc, setAudioSrc] = useState('');

  const audioChunksRef = useRef([]);

  useEffect(() => {
    if (audioChunks.length && !isRecording) {
      const audioBlob = new Blob(audioChunks, { type: 'audio/webm' });
      convertToWebM(audioBlob);
    }
  }, [audioChunks, isRecording]);

  useEffect(() => {
    if (audioUrl && audioRef.current) {
      audioRef.current.src = audioUrl;
      audioRef.current.oncanplaythrough = () => {
        audioRef.current.play().catch(err => {
          console.error('Error playing audio:', err);
        });
      };
      audioRef.current.onerror = (e) => {
        console.error('Audio playback error:', e);
      };
      audioRef.current.load();
    }
  }, [audioUrl]);

  const startRecording = async () => {
    console.log('Start recording...');
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      const options = { mimeType: 'audio/webm' };
      const recorder = new MediaRecorder(stream, options);
      setMediaRecorder(recorder);
      setAudioChunks([]);
      audioChunksRef.current = [];

      recorder.ondataavailable = event => {
        console.log('Data available:', event.data);
        if (event.data.size > 0) {
          audioChunksRef.current.push(event.data);
          setAudioChunks(prev => [...prev, event.data]);
        }
      };

      recorder.onstart = () => {
        console.log('Recording started');
      };

      recorder.onstop = async () => {
        console.log('Recording stopped');
        if (audioChunksRef.current.length === 0) {
          console.error('No audio data recorded');
          return;
        }

        const audioBlob = new Blob(audioChunksRef.current, { type: 'audio/webm' });
        console.log('Audio Blob created on stop:', audioBlob);

        const arrayBuffer = await audioBlob.arrayBuffer();
        const audioData = new Uint8Array(arrayBuffer);
        console.log('Audio Data array:', audioData);

        const transcriptionResult = await V2T(Array.from(audioData));
        console.log('Transcription result:', transcriptionResult);
        setTranscription(transcriptionResult);

        const vmtResult = await VMT(transcriptionResult);
        console.log('VMT result:', vmtResult);
        setVmtResult(vmtResult);

        const t2vResult = await T2V(vmtResult);
        console.log('T2V result length:', t2vResult.length);
        const data = await response.arrayBuffer();
        const blob = new Blob([data], { type: 'audio/mp3' });
        const url = URL.createObjectURL(blob);
        setAudioSrc(url)

       

        if (audioBlobForPlayback.size === 0) {
          console.error('Audio Blob is empty!');
          return;
        }

        const audioUrlForPlayback = URL.createObjectURL(audioBlobForPlayback);
        console.log('Audio URL for playback:', audioUrlForPlayback);
        setAudioUrl(audioUrlForPlayback);

        const audioUrlForDownload = URL.createObjectURL(audioBlobForPlayback);
        setAudioDownloadUrl(audioUrlForDownload);

        const audio = new Audio(audioUrlForPlayback);
        audio.play().catch(err => {
          console.error('Error playing audio:', err);
        });

        audio.onerror = (e) => {
          console.error('Audio playback error:', e);
        };
      };

      recorder.start();
      setIsRecording(true);
    } catch (err) {
      console.error('Error accessing media devices:', err);
    }
  };

  const stopRecording = () => {
    console.log('Stop recording...');
    if (mediaRecorder && mediaRecorder.state !== "inactive") {
      mediaRecorder.stop();
    }
    setIsRecording(false);
  };

  const convertToWebM = (audioBlob) => {
    const audioUrl = URL.createObjectURL(audioBlob);
    console.log('Recorded audio URL:', audioUrl);
    setAudioUrl(audioUrl);
  };

  return (
    <div>
      <button onClick={isRecording ? stopRecording : startRecording}>
        {isRecording ? '停止录音' : '开始录音'}
      </button>
            {audioSrc && <audio controls src={audioSrc} />}

      {audioUrl && (
        <audio ref={audioRef} controls autoPlay>
          <source src={audioUrl} type="audio/webm" />
          Your browser does not support the audio element.
        </audio>
      )}
      {transcription && (
        <div style={{ textAlign: 'right' }}>
          <h3>转录结果:</h3>
          <p>{transcription}</p>
        </div>
      )}
      {vmtResult && (
        <div style={{ textAlign: 'left' }}>
          <h3>VMT结果:</h3>
          <p>{vmtResult}</p>
        </div>
      )}
    </div>
  );
};

export default App;
