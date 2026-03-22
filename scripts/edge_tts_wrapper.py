#!/usr/bin/env python3
"""
Edge TTS Wrapper for Cortex
Generates high-quality speech using Microsoft Edge's neural voices
"""

import asyncio
import sys
import edge_tts

async def generate_speech(text: str, voice: str, output_file: str):
    """Generate speech from text using Edge TTS"""
    try:
        communicate = edge_tts.Communicate(text, voice)
        await communicate.save(output_file)
        print(f"✓ Generated audio: {output_file}", file=sys.stderr)
        return 0
    except Exception as e:
        print(f"✗ Error generating audio: {e}", file=sys.stderr)
        return 1

def main():
    if len(sys.argv) != 4:
        print("Usage: edge_tts_wrapper.py <text> <voice> <output_file>", file=sys.stderr)
        print("Example: edge_tts_wrapper.py 'Hello world' en-US-AriaNeural output.mp3", file=sys.stderr)
        sys.exit(1)

    text = sys.argv[1]
    voice = sys.argv[2]
    output_file = sys.argv[3]

    exit_code = asyncio.run(generate_speech(text, voice, output_file))
    sys.exit(exit_code)

if __name__ == "__main__":
    main()
