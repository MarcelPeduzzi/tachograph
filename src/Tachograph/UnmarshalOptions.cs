using System;
using System.IO;

namespace Tachograph
{
    /// <summary>
    /// Configures the unmarshaling process for tachograph files.
    /// </summary>
    public class UnmarshalOptions
    {
        /// <summary>
        /// Controls how the unmarshaler handles unrecognized tags or structural inconsistencies.
        /// 
        /// If true (default), the unmarshaler will return an error on any unrecognized tags or structural inconsistencies.
        /// If false, the unmarshaler will attempt to skip over unrecognized parts of the file and continue parsing.
        /// </summary>
        public bool Strict { get; set; } = true;

        /// <summary>
        /// Parses a tachograph file from its binary representation into a raw, unparsed format.
        /// The returned RawFile is suitable for authentication.
        /// </summary>
        /// <param name="data">The binary data to unmarshal</param>
        /// <returns>A RawFile object</returns>
        /// <exception cref="InvalidDataException">Thrown when the data is invalid or unrecognized</exception>
        public RawFile Unmarshal(byte[] data)
        {
            if (data == null)
                throw new ArgumentNullException(nameof(data));

            if (data.Length < 2)
                throw new InvalidDataException("Insufficient data for tachograph file");

            // Vehicle unit file (starts with TREP prefix 0x76)
            if (data[0] == 0x76)
            {
                var vuRaw = UnmarshalVehicleUnitFile(data);
                return new RawFile
                {
                    Type = RawFileType.VehicleUnit,
                    VehicleUnit = vuRaw
                };
            }

            // Card file (starts with EF_ICC prefix 0x0002)
            if (data.Length >= 2)
            {
                ushort prefix = (ushort)((data[0] << 8) | data[1]);
                if (prefix == 0x0002)
                {
                    var cardRaw = UnmarshalCardFile(data);
                    return new RawFile
                    {
                        Type = RawFileType.Card,
                        Card = cardRaw
                    };
                }
            }

            throw new InvalidDataException("Unknown or unsupported file type");
        }

        private RawCardFile UnmarshalCardFile(byte[] data)
        {
            // TODO: Implement card file unmarshaling
            throw new NotImplementedException("Card file unmarshaling not yet implemented");
        }

        private RawVehicleUnitFile UnmarshalVehicleUnitFile(byte[] data)
        {
            // TODO: Implement VU file unmarshaling
            throw new NotImplementedException("Vehicle unit file unmarshaling not yet implemented");
        }
    }
}
