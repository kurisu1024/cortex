#!/usr/bin/env python3
"""
Stable Diffusion Image Generator for Cortex
Generates AI images from text prompts using Stable Diffusion
"""

import sys
import torch
from diffusers import StableDiffusionPipeline
import os
import numpy as np

def generate_image(prompt: str, output_file: str, model_id: str = "runwayml/stable-diffusion-v1-5"):
    """Generate an image from a text prompt using Stable Diffusion"""
    try:
        print(f"Loading Stable Diffusion model: {model_id}", file=sys.stderr)

        # MPS (Apple Silicon) has known issues with Stable Diffusion causing black images
        # Use CPU for reliability - it's slower but produces correct images
        # See: https://github.com/pytorch/pytorch/issues/77764
        device = "cpu"
        dtype = torch.float32
        print("Using CPU (most reliable for Stable Diffusion)", file=sys.stderr)
        print("Note: Image generation will take ~2-3 minutes per image", file=sys.stderr)

        # Load the pipeline
        pipe = StableDiffusionPipeline.from_pretrained(
            model_id,
            torch_dtype=dtype,
            safety_checker=None,  # Disable safety checker for speed
        )
        pipe = pipe.to(device)

        # Enable attention slicing to reduce memory usage
        pipe.enable_attention_slicing()

        print(f"Generating image for prompt: {prompt[:50]}...", file=sys.stderr)

        # Generate the image
        with torch.no_grad():
            result = pipe(
                prompt,
                num_inference_steps=30,  # Reduce steps for faster generation
                guidance_scale=7.5,
                height=768,  # 16:9 aspect ratio
                width=1360
            )
            image = result.images[0]

        # Validate the image isn't corrupted (all black)
        img_array = np.array(image)
        mean_value = np.mean(img_array)
        std_value = np.std(img_array)

        if mean_value < 10 or std_value < 5:
            print(f"✗ Error: Generated image is corrupted (all black)", file=sys.stderr)
            print(f"  Image stats: mean={mean_value:.2f}, std={std_value:.2f}", file=sys.stderr)
            return 1

        # Save the image
        image.save(output_file)
        print(f"✓ Generated image: {output_file}", file=sys.stderr)
        print(f"  Image quality: mean={mean_value:.2f}, std={std_value:.2f}", file=sys.stderr)
        return 0

    except Exception as e:
        print(f"✗ Error generating image: {e}", file=sys.stderr)
        import traceback
        traceback.print_exc(file=sys.stderr)
        return 1

def main():
    if len(sys.argv) < 3:
        print("Usage: sd_image_gen.py <prompt> <output_file> [model_id]", file=sys.stderr)
        print("Example: sd_image_gen.py 'A beautiful sunset' output.png", file=sys.stderr)
        print("Optional model_id: runwayml/stable-diffusion-v1-5 (default)", file=sys.stderr)
        sys.exit(1)

    prompt = sys.argv[1]
    output_file = sys.argv[2]
    model_id = sys.argv[3] if len(sys.argv) > 3 else "runwayml/stable-diffusion-v1-5"

    exit_code = generate_image(prompt, output_file, model_id)
    sys.exit(exit_code)

if __name__ == "__main__":
    main()
