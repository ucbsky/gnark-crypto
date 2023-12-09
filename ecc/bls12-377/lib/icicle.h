#include <stdbool.h>

int msm_cuda_bls12_377(
  void* out, void* points, void* scalars, size_t count, unsigned large_bucket_factor, size_t device_id);
