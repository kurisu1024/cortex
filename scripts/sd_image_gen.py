#!/usr/bin/env python3
"""
Stable Diffusion Image Generator for Cortex
Generates AI images from text prompts using Stable Diffusion
Supports SDXL-Turbo for 10-20x faster generation on Apple Silicon
"""

import sys
import torch
from diffusers import StableDiffusionPipeline, AutoPipelineForText2Image
import os
import numpy as np

# Model presets for different speed/quality tradeoffs
MODEL_PRESETS = {
    "sd15": "runwayml/stable-diffusion-v1-5",  # Default, good quality, slower (30 steps)
    "sdxl-turbo": "stabilityai/sdxl-turbo",     # Fastest, 1 step, excellent quality
    "sd-turbo": "stabilityai/sd-turbo",         # Fast, 1 step, good quality, lower res
}

def generate_image(prompt: str, output_file: str, model_id: str = "runwayml/stable-diffusion-v1-5"):
    """Generate an image from a text prompt using Stable Diffusion"""
    # Resolve model preset if using shorthand
    if model_id in MODEL_PRESETS:
        model_id = MODEL_PRESETS[model_id]

    is_turbo = "turbo" in model_id.lower()

    # Try MPS (GPU) first for turbo models, fall back to CPU if it fails
    devices_to_try = []

    if torch.backends.mps.is_available() and is_turbo:
        # Turbo models work better on MPS
        devices_to_try.append(("mps", torch.float32, "MPS (Apple Silicon GPU)"))

    devices_to_try.append(("cpu", torch.float32, "CPU"))

    last_error = None

    for device, dtype, device_name in devices_to_try:
        try:
            print(f"Loading model: {model_id}", file=sys.stderr)
            print(f"Using {device_name} for generation", file=sys.stderr)

            if device == "mps":
                if is_turbo:
                    print("Note: ~5-10 seconds per image on M4 (Turbo mode)", file=sys.stderr)
                else:
                    print("Note: ~30-60 seconds per image on M4", file=sys.stderr)
                os.environ["PYTORCH_ENABLE_MPS_FALLBACK"] = "1"
            else:
                if is_turbo:
                    print("Note: ~20-30 seconds per image (Turbo mode)", file=sys.stderr)
                else:
                    print("Note: Image generation will take ~2-3 minutes per image", file=sys.stderr)

            # Load the pipeline - use AutoPipeline for turbo models
            if is_turbo:
                pipe = AutoPipelineForText2Image.from_pretrained(
                    model_id,
                    torch_dtype=dtype,
                    variant="fp16" if device == "mps" else None,
                )
            else:
                pipe = StableDiffusionPipeline.from_pretrained(
                    model_id,
                    torch_dtype=dtype,
                    safety_checker=None,
                    variant="fp16" if device == "mps" else None,
                )

            pipe = pipe.to(device)

            # Enable memory optimizations
            pipe.enable_attention_slicing()
            if device == "mps":
                pipe.enable_vae_slicing()

            print(f"Generating image for prompt: {prompt[:50]}...", file=sys.stderr)

            # Set parameters based on model type
            if is_turbo:
                # Turbo models work best with 1 step and no guidance
                num_steps = 1
                guidance = 0.0
                height = 512
                width = 512
            else:
                num_steps = 30
                guidance = 7.5
                height = 768
                width = 1360

            # Generate the image
            with torch.no_grad():
                result = pipe(
                    prompt,
                    num_inference_steps=num_steps,
                    guidance_scale=guidance,
                    height=height,
                    width=width,
                )
                image = result.images[0]

            # Validate the image isn't corrupted (all black)
            img_array = np.array(image)
            mean_value = np.mean(img_array)
            std_value = np.std(img_array)

            if mean_value < 10 or std_value < 5:
                raise ValueError(f"Generated image is corrupted (all black): mean={mean_value:.2f}, std={std_value:.2f}")

            # Save the image
            image.save(output_file)
            print(f"✓ Generated image: {output_file}", file=sys.stderr)
            print(f"  Image quality: mean={mean_value:.2f}, std={std_value:.2f}", file=sys.stderr)
            return 0

        except Exception as e:
            last_error = e
            print(f"✗ Error with {device_name}: {e}", file=sys.stderr)

            # If there are more devices to try, continue
            if device != devices_to_try[-1][0]:
                print(f"Falling back to next device...", file=sys.stderr)
                continue
            else:
                # This was the last device, give up
                print(f"✗ All devices failed", file=sys.stderr)
                import traceback
                traceback.print_exc(file=sys.stderr)
                return 1

    # Should never reach here
    print(f"✗ Unexpected error: {last_error}", file=sys.stderr)
    return 1

    # Should never reach here
    print(f"✗ Unexpected error: {last_error}", file=sys.stderr)
    return 1

def main():
    if len(sys.argv) < 3:
        print("Usage: sd_image_gen.py <prompt> <output_file> [model_id]", file=sys.stderr)
        print("Example: sd_image_gen.py 'A beautiful sunset' output.png sdxl-turbo", file=sys.stderr)
        print("", file=sys.stderr)
        print("Model options:", file=sys.stderr)
        print("  sdxl-turbo  - Fastest (5-10s on M4 GPU), best quality", file=sys.stderr)
        print("  sd-turbo    - Fast (3-5s on M4 GPU), good quality", file=sys.stderr)
        print("  sd15        - Slower (2-3min), good quality", file=sys.stderr)
        print("  Or use full model ID like 'stabilityai/sdxl-turbo'", file=sys.stderr)
        sys.exit(1)

    prompt = sys.argv[1]
    output_file = sys.argv[2]
    model_id = sys.argv[3] if len(sys.argv) > 3 else "stabilityai/sdxl-turbo"

    exit_code = generate_image(prompt, output_file, model_id)
    sys.exit(exit_code)

if __name__ == "__main__":
    main()
