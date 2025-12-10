using System;
using System.Threading;
using System.Threading.Tasks;

namespace Tachograph
{
    /// <summary>
    /// Main API for working with tachograph files.
    /// Provides methods for unmarshaling, parsing, authenticating, anonymizing, and marshaling tachograph data.
    /// </summary>
    public static class Tachograph
    {
        /// <summary>
        /// Parses a tachograph file from its binary representation into a raw, unparsed format with default options.
        /// The returned RawFile is suitable for authentication.
        /// </summary>
        /// <param name="data">The binary data to unmarshal</param>
        /// <returns>A RawFile object</returns>
        public static RawFile Unmarshal(byte[] data)
        {
            var opts = new UnmarshalOptions
            {
                Strict = true
            };
            return opts.Unmarshal(data);
        }

        /// <summary>
        /// Performs semantic parsing on raw tachograph records with default options.
        /// If the raw file has been authenticated (via Authenticate), the authentication results are propagated to the parsed messages.
        /// </summary>
        /// <param name="rawFile">The raw file to parse</param>
        /// <returns>A parsed File object</returns>
        public static File Parse(RawFile rawFile)
        {
            var opts = new ParseOptions
            {
                PreserveRawData = true
            };
            return opts.Parse(rawFile);
        }

        /// <summary>
        /// Performs cryptographic authentication on a raw file with default options.
        /// </summary>
        /// <param name="rawFile">The raw file to authenticate</param>
        /// <param name="cancellationToken">Cancellation token</param>
        /// <returns>An authenticated RawFile</returns>
        public static async Task<RawFile> AuthenticateAsync(RawFile rawFile, CancellationToken cancellationToken = default)
        {
            var opts = new AuthenticateOptions
            {
                Mutate = false
            };
            return await opts.AuthenticateAsync(rawFile, cancellationToken);
        }

        /// <summary>
        /// Creates an anonymized copy of a parsed tachograph file with default options.
        /// </summary>
        /// <param name="file">The file to anonymize</param>
        /// <returns>An anonymized File object</returns>
        public static File Anonymize(File file)
        {
            return new AnonymizeOptions().Anonymize(file);
        }

        /// <summary>
        /// Serializes a parsed tachograph file into binary format with default options.
        /// </summary>
        /// <param name="file">The file to marshal</param>
        /// <returns>Binary data</returns>
        public static byte[] Marshal(File file)
        {
            var opts = new MarshalOptions
            {
                UseRawData = true
            };
            return opts.Marshal(file);
        }
    }
}
