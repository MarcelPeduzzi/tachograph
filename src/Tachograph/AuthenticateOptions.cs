using System;
using System.Threading;
using System.Threading.Tasks;

namespace Tachograph
{
    /// <summary>
    /// Configures the signature authentication process.
    /// </summary>
    public class AuthenticateOptions
    {
        /// <summary>
        /// Used to resolve CA certificates by their Certificate Authority Reference (CAR).
        /// If null, this defaults to using DefaultCertificateResolver.
        /// </summary>
        public ICertificateResolver CertificateResolver { get; set; }

        /// <summary>
        /// Controls whether authentication modifies the input RawFile in-place.
        /// 
        /// If false (default), the input RawFile is deep cloned before authentication,
        /// ensuring the original remains unchanged. This is the safe default for most use cases.
        /// 
        /// If true, the input RawFile is modified in-place, which is more efficient
        /// but requires the caller to be aware that the input will be mutated.
        /// </summary>
        public bool Mutate { get; set; } = false;

        /// <summary>
        /// Performs cryptographic authentication on a raw tachograph file,
        /// populating Authentication fields in the raw records.
        /// 
        /// By default (Mutate: false), this method returns a new authenticated RawFile,
        /// leaving the input unchanged. Set Mutate: true for in-place authentication.
        /// </summary>
        /// <param name="rawFile">The raw file to authenticate</param>
        /// <param name="cancellationToken">Cancellation token</param>
        /// <returns>An authenticated RawFile</returns>
        /// <exception cref="ArgumentNullException">Thrown when rawFile is null</exception>
        public async Task<RawFile> AuthenticateAsync(RawFile rawFile, CancellationToken cancellationToken = default)
        {
            if (rawFile == null)
                throw new ArgumentNullException(nameof(rawFile));

            // Clone the input unless mutate is explicitly requested
            var target = Mutate ? rawFile : rawFile.Clone();

            var resolver = CertificateResolver ?? DefaultCertificateResolver.Instance;

            switch (target.Type)
            {
                case RawFileType.Card:
                    await AuthenticateCardFileAsync(target.Card, resolver, cancellationToken);
                    break;

                case RawFileType.VehicleUnit:
                    await AuthenticateVehicleUnitFileAsync(target.VehicleUnit, resolver, cancellationToken);
                    break;

                default:
                    throw new NotSupportedException($"Unknown raw file type: {target.Type}");
            }

            return target;
        }

        private async Task AuthenticateCardFileAsync(RawCardFile card, ICertificateResolver resolver, CancellationToken cancellationToken)
        {
            // TODO: Implement card file authentication
            await Task.CompletedTask;
            throw new NotImplementedException("Card file authentication not yet implemented");
        }

        private async Task AuthenticateVehicleUnitFileAsync(RawVehicleUnitFile vu, ICertificateResolver resolver, CancellationToken cancellationToken)
        {
            // TODO: Implement VU file authentication
            await Task.CompletedTask;
            throw new NotImplementedException("Vehicle unit file authentication not yet implemented");
        }
    }
}
