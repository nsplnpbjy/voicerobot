import React, { useState, useEffect, useRef } from 'react';
import { V2T, T2V, VMT } from "../wailsjs/go/main/App";
import { Button, Container, CssBaseline, AppBar, Toolbar, Typography, Paper, Box, IconButton } from '@mui/material';
import { ThemeProvider, createTheme, responsiveFontSizes } from '@mui/material/styles';
import PauseIcon from '@mui/icons-material/Pause';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import MicIcon from '@mui/icons-material/Mic';
import MicOffIcon from '@mui/icons-material/MicOff';
import './App.css'; // 导入自定义CSS文件

let theme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#bb86fc',
    },
    secondary: {
      main: '#03dac6',
    },
    background: {
      default: '#121212',
      paper: '#1d1d1d',
    },
    text: {
      primary: '#ffffff',
      secondary: '#a1a1a1',
    },
  },
});

theme = responsiveFontSizes(theme);

const App = () => {
  const [isRecording, setIsRecording] = useState(false);
  const [audioChunks, setAudioChunks] = useState([]);
  const [mediaRecorder, setMediaRecorder] = useState(null);
  const [audioSrc, setAudioSrc] = useState('');
  const [messages, setMessages] = useState([]);
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentAudio, setCurrentAudio] = useState(null);

  const audioChunksRef = useRef([]);
  const audioRefs = useRef([]);
  const messagesEndRef = useRef(null);

  useEffect(() => {
    if (audioChunks.length && !isRecording) {
      const audioBlob = new Blob(audioChunks, { type: 'audio/webm' });
      convertToWebM(audioBlob);
    }
  }, [audioChunks, isRecording]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  const startRecording = async () => {
    try {
      stopAllAudios();
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      const recorder = new MediaRecorder(stream);
      setMediaRecorder(recorder);
      setAudioChunks([]);
      audioChunksRef.current = [];
      recorder.ondataavailable = event => {
        if (event.data.size > 0) {
          audioChunksRef.current.push(event.data);
          setAudioChunks(prev => [...prev, event.data]);
        }
      };
      recorder.onstop = async () => {
        if (audioChunksRef.current.length === 0) {
          return;
        }
        const audioBlob = new Blob(audioChunksRef.current, { type: 'audio/webm' });
        const arrayBuffer = await audioBlob.arrayBuffer();
        const audioData = new Uint8Array(arrayBuffer);
        const transcriptionResult = await V2T(Array.from(audioData));
        setMessages(prev => [...prev, { sender: '用户', text: transcriptionResult }]);
        const vmtResult = await VMT(transcriptionResult);
        setMessages(prev => [...prev, { sender: 'AI', text: vmtResult }]);
        const t2vResult = await T2V(vmtResult);
        const byteArray = Uint8Array.from(atob(t2vResult), c => c.charCodeAt(0));
        const blob = new Blob([byteArray], { type: 'audio/mpeg' });
        const url = URL.createObjectURL(blob);
        setAudioSrc(url);
        const audio = new Audio(url);
        audio.onloadeddata = () => {
          stopAllAudios();
          audio.play();
          setIsPlaying(true);
          setCurrentAudio(audio);

          audioRefs.current.push(audio);
          audio.onpause = () => setIsPlaying(false);
          audio.onplay = () => setIsPlaying(true);
        };
        audio.load();
      };
      recorder.start();
      setIsRecording(true);
    } catch (err) {
      console.error('Error accessing media devices:', err);
    }
  };

  const stopRecording = () => {
    if (mediaRecorder && mediaRecorder.state !== "inactive") {
      mediaRecorder.stop();
    }
    setIsRecording(false);
  };

  const convertToWebM = (audioBlob) => {
    const audioUrl = URL.createObjectURL(audioBlob);
    console.log('Recorded audio URL:', audioUrl);
  };

  const stopAllAudios = () => {
    audioRefs.current.forEach(audio => {
      if (!audio.paused) {
        audio.pause();
        audio.currentTime = 0;
      }
    });
    setCurrentAudio(null);
  };

  const handlePlayPause = () => {
    if (isPlaying) {
      currentAudio.pause();
      setIsPlaying(false);
    } else {
      if (currentAudio) {
        currentAudio.play();
        setIsPlaying(true);
      }
    }
  };

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6">白胖的语音AI</Typography>
        </Toolbar>
      </AppBar>
      <Container maxWidth="md" style={{ marginTop: '20px', height: 'calc(100% - 64px)', position: 'relative', backgroundColor: '#121212' }}>
        <Box
          position="fixed"
          bottom={20}
          left={20}
          zIndex="tooltip"
          display="flex"
          flexDirection="column"
          gap={2}
        >
          <IconButton
            onClick={isRecording ? stopRecording : startRecording}
            color={isRecording ? "secondary" : "primary"}
            className={isRecording ? "recording" : ""}
          >
            {isRecording ? <MicIcon /> : <MicOffIcon />}
          </IconButton>
          <IconButton onClick={handlePlayPause} color="primary">
            {isPlaying ? <PauseIcon /> : <PlayArrowIcon />}
          </IconButton>
        </Box>
        <Box mt={3}>
          {messages.map((msg, index) => (
            <Paper key={index} style={{ padding: '10px', marginBottom: '10px', textAlign: msg.sender === '用户' ? 'right' : 'left', backgroundColor: '#1d1d1d' }}>
              <Typography variant="body1"><strong>{msg.sender}:</strong> {msg.text}</Typography>
            </Paper>
          ))}
          <div ref={messagesEndRef} />
        </Box>
        {audioSrc && (
          <Box mt={3}>
            <audio controls style={{ display: 'none' }} src={audioSrc}>
              Your browser does not support the audio element.
            </audio>
          </Box>
        )}
      </Container>
    </ThemeProvider>
  );
};

export default App;
